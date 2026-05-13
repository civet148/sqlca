package sqlca

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/civet148/sqlca/v3/types"
)

func (e *Engine) addPreload(query string, args ...any) *Engine {
	if e.preloads == nil {
		e.preloads = make(map[string][]any)
	}
	e.preloads[query] = args
	return e
}

func (e *Engine) handlePreloads() (err error) {
	for query, args := range e.preloads {
		if err = e.execPreload(query, args...); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) execPreload(query string, args ...any) (err error) {
	if e.model == nil {
		return fmt.Errorf("model is nil, please call Model method first")
	}

	// 获取当前模型类型和值
	modelType := reflect.TypeOf(e.model)
	modelValue := reflect.ValueOf(e.model)

	// 处理指针类型
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}

	// 如果是切片类型，使用批量加载模式（解决N+1问题）
	if modelType.Kind() == reflect.Slice {
		return e.execPreloadSlice(modelValue, query, args...)
	}

	// 单个结构体实例
	return e.loadAssociations(modelValue, query, args...)
}

// execPreloadSlice 批量预加载关联数据（解决N+1问题）
func (e *Engine) execPreloadSlice(modelValue reflect.Value, associationName string, args ...any) error {
	n := modelValue.Len()
	if n == 0 {
		return nil
	}

	// 提取所有元素的反射值（解引用指针）
	elements := make([]reflect.Value, n)
	for i := 0; i < n; i++ {
		element := modelValue.Index(i)
		if element.Kind() == reflect.Ptr {
			element = element.Elem()
		}
		elements[i] = element
	}

	// 收集所有元素的主键值
	pkValues := make([]interface{}, n)
	allValid := true
	for i, element := range elements {
		pkValues[i] = getPKValue(element)
		if pkValues[i] == nil {
			allValid = false
		}
	}

	if !allValid {
		// 如果有元素没有主键值，降级为逐个加载
		for i, element := range elements {
			if err := e.loadAssociations(element, associationName, args...); err != nil {
				return err
			}
			_ = i
		}
		return nil
	}

	return e.loadAssociationsBatch(elements, pkValues, associationName, args...)
}

// loadAssociationsBatch 批量加载关联数据（按字段索引分发）
func (e *Engine) loadAssociationsBatch(elements []reflect.Value, pkValues []interface{}, associationName string, args ...any) error {
	if len(elements) == 0 {
		return nil
	}

	modelType := elements[0].Type()

	// 遍历结构体字段寻找匹配的关联字段
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// 检查字段名是否匹配关联名称
		if field.Name != associationName {
			continue
		}

		// 获取GORM标签
		gormTag := recursiveTags(field, types.TAG_NAME_GORM)
		if gormTag == "" {
			continue
		}

		// 解析GORM标签设置
		settings := parseTagSetting(gormTag, ";")

		// 检查是否为many2many关联
		if junctionTable, ok := settings["MANY2MANY"]; ok {
			return e.loadMany2ManyAssociationBatch(elements, pkValues, i, field.Type, junctionTable, args...)
		}

		// 检查是否为foreignKey关联
		if foreignKey, ok := settings["FOREIGNKEY"]; ok {
			return e.loadForeignKeyAssociationBatch(elements, pkValues, i, field.Type, foreignKey, args...)
		}

		// 检查是否有references引用
		if references, ok := settings["REFERENCES"]; ok {
			foreignKey, _ := settings["FOREIGNKEY"]
			return e.loadReferenceAssociationBatch(elements, pkValues, i, field.Type, foreignKey, references, args...)
		}
	}

	return fmt.Errorf("association '%s' not found", associationName)
}

// loadAssociations 加载关联数据
func (e *Engine) loadAssociations(modelValue reflect.Value, associationName string, args ...any) error {
	modelType := modelValue.Type()

	// 遍历结构体字段寻找匹配的关联字段
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		// 检查字段名是否匹配关联名称
		if field.Name != associationName {
			continue
		}

		// 获取GORM标签
		gormTag := recursiveTags(field, types.TAG_NAME_GORM, types.TAG_NAME_SQLCA)
		if gormTag == "" {
			continue
		}

		// 解析GORM标签设置
		settings := parseTagSetting(gormTag, ";")

		// 检查是否为many2many关联
		if junctionTable, ok := settings["MANY2MANY"]; ok {
			return e.loadMany2ManyAssociation(modelValue, fieldValue, field.Type, junctionTable, args...)
		}

		// 检查是否为foreignKey关联
		if foreignKey, ok := settings["FOREIGNKEY"]; ok {
			return e.loadForeignKeyAssociation(modelValue, fieldValue, field.Type, foreignKey, args...)
		}

		// 检查是否有references引用
		if references, ok := settings["REFERENCES"]; ok {
			foreignKey, _ := settings["FOREIGNKEY"]
			return e.loadReferenceAssociation(modelValue, fieldValue, field.Type, foreignKey, references, args...)
		}
	}

	return fmt.Errorf("association '%s' not found", associationName)
}

// loadMany2ManyAssociation 加载多对多关联
func (e *Engine) loadMany2ManyAssociation(mainModelValue reflect.Value, fieldValue reflect.Value, fieldType reflect.Type, junctionTable string, args ...any) error {
	if fieldValue.Kind() != reflect.Ptr && fieldValue.Kind() != reflect.Slice {
		return fmt.Errorf("many2many association field must be a pointer to slice")
	}

	// 如果是指针，获取其指向的切片
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			newSlice := reflect.New(fieldType.Elem())
			fieldValue.Set(newSlice)
			fieldValue = newSlice.Elem()
		} else {
			fieldValue = fieldValue.Elem()
		}
	}

	if fieldValue.Kind() != reflect.Slice {
		return fmt.Errorf("many2many association field must be a pointer to slice")
	}

	// 处理主模型指针
	mainModelReflectValue := mainModelValue
	if mainModelReflectValue.Kind() == reflect.Ptr {
		mainModelReflectValue = mainModelReflectValue.Elem()
	}

	mainModelType := mainModelReflectValue.Type()

	// 尝试找到主键字段 (通常名为ID)
	var mainPKValue interface{}
	for i := 0; i < mainModelType.NumField(); i++ {
		field := mainModelType.Field(i)
		gormTag := recursiveTags(field, types.TAG_NAME_GORM, types.TAG_NAME_SQLCA)
		settings := parseTagSetting(gormTag, ";")
		if _, isPK := settings["PRIMARYKEY"]; isPK || strings.EqualFold(field.Name, "Id") || strings.EqualFold(field.Name, "ID") {
			fieldVal := mainModelReflectValue.Field(i)
			mainPKValue = fieldVal.Interface()
			break
		}
	}

	if mainPKValue == nil {
		// 如果没找到明确的主键，尝试获取第一个整数类型的字段作为主键
		for i := 0; i < mainModelType.NumField(); i++ {
			field := mainModelType.Field(i)
			fieldVal := mainModelReflectValue.Field(i)
			if strings.EqualFold(field.Name, "Id") || strings.EqualFold(field.Name, "ID") {
				mainPKValue = fieldVal.Interface()
				break
			}
		}
	}

	if mainPKValue == nil {
		return fmt.Errorf("could not find primary key value for main model")
	}

	// 构建关联查询：从junctionTable中查找对应记录，再关联到目标表
	elementType := fieldValue.Type().Elem()
	if elementType.Kind() == reflect.Ptr {
		elementType = elementType.Elem()
	}

	// 创建关联模型实例
	_ = reflect.New(elementType).Interface()

	// 执行查询：SELECT * FROM assoc_table WHERE id IN (SELECT assoc_id FROM junction_table WHERE main_id = ?)
	// 首先从junctionTable获取关联ID列表
	// 解析主模型和关联模型的字段名
	mainModelName := getTableNameFromType(mainModelType)
	assocModelName := getTableNameFromType(elementType)

	// 生成可能的外键列名
	junctionPK1 := mainModelName + "_id"
	junctionPK2 := assocModelName + "_id"

	// 尝试从GORM标签中获取更准确的外键名称
	// 如果主模型是User，关联模型是Role，则junctionPK1=user_id, junctionPK2=role_id
	junctionPK1 = snakeCase(mainModelType.Name()) + "_id"
	junctionPK2 = snakeCase(elementType.Name()) + "_id"

	var assocIDs []interface{}
	junctionQuery := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", junctionPK2, junctionTable, junctionPK1)
	rows, err := e.db.Query(junctionQuery, mainPKValue)
	if err != nil {
		return fmt.Errorf("failed to query junction table: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var assocID interface{}
		if err := rows.Scan(&assocID); err != nil {
			return fmt.Errorf("failed to scan junction row: %v", err)
		}
		assocIDs = append(assocIDs, assocID)
	}

	if len(assocIDs) == 0 {
		// 没有关联的数据，设置空切片
		fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
		return nil
	}

	// 使用IN查询获取关联数据
	// 合并关联ID和额外的查询条件
	totalArgs := make([]interface{}, 0, len(assocIDs)+len(args))
	totalArgs = append(totalArgs, assocIDs...)
	totalArgs = append(totalArgs, args...)

	placeholders := make([]string, len(assocIDs))
	queryArgs := make([]interface{}, len(assocIDs))
	for i, id := range assocIDs {
		placeholders[i] = "?"
		queryArgs[i] = id
	}

	tableName := getTableNameFromType(elementType)

	// 动态构建选择字段列表，而不是使用 *
	columns := getStructFields(elementType, e)
	selectFields := strings.Join(columns, ", ")

	inClause := strings.Join(placeholders, ", ")

	// 构建基础查询
	finalQuery := fmt.Sprintf("SELECT %s FROM %s WHERE id IN (%s)", selectFields, tableName, inClause)

	// 如果有额外的查询条件，则添加它们
	if len(args) > 0 {
		// 支持GORM风格的条件参数，例如: "status = ? AND deleted_at is NULL", 1
		// 或者复合条件: "status = ? AND created_at > ?", "active", "2023-01-01"
		// 检查第一个参数是否为字符串（可能是条件表达式）
		if condition, ok := args[0].(string); ok {
			// 添加条件到查询中
			finalQuery += " AND " + condition
			// 将其余参数添加到查询参数中
			if len(args) > 1 {
				queryArgs = append(queryArgs, args[1:]...)
			}
		}
	}

	// 替换占位符为正确的数据库方言格式
	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		// 重新构建查询以适应PostgreSQL占位符
		// 为IN子句中的ID创建占位符
		inPlace := make([]string, len(assocIDs))
		for i := 0; i < len(assocIDs); i++ {
			inPlace[i] = "$" + strconv.Itoa(i+1)
		}
		inClause = strings.Join(inPlace, ", ")

		// 重新构建基础查询
		finalQuery = fmt.Sprintf("SELECT %s FROM %s WHERE id IN (%s)", selectFields, tableName, inClause)

		// 添加额外条件
		argOffset := len(assocIDs) + 1
		if len(args) > 0 {
			// 检查第一个参数是否为字符串（可能是条件表达式）
			if condition, ok := args[0].(string); ok {
				// 处理包含占位符的字符串条件
				// 计算条件字符串中占位符的数量，并调整argOffset
				placeholderCount := strings.Count(condition, "?")
				finalQuery += " AND " + strings.ReplaceAll(condition, "?", fmt.Sprintf("$%d", argOffset))

				// 更新argOffset为占位符之后的位置
				argOffset += placeholderCount
			} else if len(args)%2 == 0 {
				for i := 0; i < len(args); i += 2 {
					finalQuery += " AND "

					field, ok := args[i].(string)
					if !ok {
						return fmt.Errorf("invalid field name type: %T", args[i])
					}

					operator := "="
					valueIndex := i + 1

					// 检查下一个参数是否是操作符
					if valueIndex < len(args) {
						nextArg, isOp := args[valueIndex].(string)
						if isOp && isInOperator(nextArg) {
							operator = nextArg
							valueIndex++ // 如果下一个是操作符，则实际值是再下一个
						}
					}

					finalQuery += fmt.Sprintf("%s %s $%d", field, operator, argOffset)
					argOffset++
				}
			}
		}

		// 为PostgreSQL重新构建参数列表
		queryArgs = make([]interface{}, 0, len(assocIDs)+len(args))
		queryArgs = append(queryArgs, assocIDs...)
		if len(args) > 0 {
			// 检查第一个参数是否为字符串（可能是条件表达式）
			if _, ok := args[0].(string); ok {
				// 如果是字符串条件表达式，直接将剩余参数添加到查询参数中
				if len(args) > 1 {
					queryArgs = append(queryArgs, args[1:]...)
				}
			} else {
				// 如果不是字符串条件表达式，按照键值对方式处理
				if len(args)%2 == 0 {
					for i := 0; i < len(args); i += 2 {
						// 跳过操作符，只添加值
						valueIndex := i + 1
						// 检查下一个参数是否是操作符
						if valueIndex < len(args) {
							nextArg, isOp := args[valueIndex].(string)
							if isOp && isInOperator(nextArg) {
								valueIndex++ // 如果下一个是操作符，则实际值是再下一个
							}
						}
						if valueIndex < len(args) {
							queryArgs = append(queryArgs, args[valueIndex])
						}
					}
				} else {
					// 如果参数不是偶数个，返回错误
					return fmt.Errorf("unsupported argument format for additional conditions")
				}
			}
		}

	}

	// 创建关联模型切片
	sliceType := reflect.SliceOf(reflect.PtrTo(elementType))
	resultSlice := reflect.New(sliceType).Elem()

	// 临时设置模型为关联模型切片
	originalModel := e.model
	e.model = resultSlice.Addr().Interface()

	// 查询关联数据
	rows, err = e.db.Query(finalQuery, queryArgs...)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to query associated data: %v", err)
	}
	defer rows.Close()

	// 使用fetchRows填充数据
	_, err = e.fetchRows(rows)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to fetch associated data: %v", err)
	}

	// 恢复原始模型
	e.model = originalModel

	// 设置关联字段值
	fieldValue.Set(resultSlice)

	return nil
}

// loadForeignKeyAssociation 加载外键关联
func (e *Engine) loadForeignKeyAssociation(mainModelValue reflect.Value, fieldValue reflect.Value, fieldType reflect.Type, foreignKey string, args ...any) error {
	// 处理主模型指针
	mainModelReflectValue := mainModelValue
	if mainModelReflectValue.Kind() == reflect.Ptr {
		mainModelReflectValue = mainModelReflectValue.Elem()
	}

	mainModelType := mainModelReflectValue.Type()

	// 查找主模型的主键字段值
	var fkValue interface{}
	found := false
	for i := 0; i < mainModelType.NumField(); i++ {
		field := mainModelType.Field(i)
		gormTag := recursiveTags(field, types.TAG_NAME_GORM, types.TAG_NAME_SQLCA)
		settings := parseTagSetting(gormTag, ";")
		_, isPK := settings["PRIMARYKEY"]
		if isPK || strings.EqualFold(field.Name, "Id") || strings.EqualFold(field.Name, "ID") {
			fieldVal := mainModelReflectValue.Field(i)
			fkValue = fieldVal.Interface()
			found = true
			break
		}
	}

	if !found {
		// 如果没找到明确的主键，尝试获取第一个整数类型的字段作为主键
		for i := 0; i < mainModelType.NumField(); i++ {
			field := mainModelType.Field(i)
			fieldVal := mainModelReflectValue.Field(i)
			if strings.EqualFold(field.Name, "Id") || strings.EqualFold(field.Name, "ID") {
				fkValue = fieldVal.Interface()
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("could not find primary key value for main model")
	}

	// 准备关联模型
	var assocModel interface{}
	if fieldType.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			// 创建新的指针实例
			newPtr := reflect.New(fieldType.Elem())
			fieldValue.Set(newPtr)
			assocModel = newPtr.Interface()
		} else {
			assocModel = fieldValue.Interface()
		}
	} else {
		// 对于非指针类型，需要创建一个指针指向它
		assocModel = fieldValue.Addr().Interface()
	}

	// 构建查询：SELECT {fields} FROM assoc_table WHERE id = fk_value
	elementType := fieldType
	if fieldType.Kind() == reflect.Ptr {
		elementType = fieldType.Elem()
	}

	tableName := getTableNameFromType(elementType)

	// 将 foreignKey 转换为蛇形命名法，以匹配数据库中的列名
	foreignKeySnake := snakeCase(foreignKey)

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, foreignKeySnake)

	// 处理不同数据库的占位符
	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tableName, foreignKeySnake)
	}

	// 临时设置模型为关联模型
	originalModel := e.model
	e.model = assocModel

	// 查询关联数据
	rows, err := e.db.Query(query, fkValue)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to query foreign key association: %v", err)
	}
	defer rows.Close()

	// 使用fetchRows填充数据
	_, err = e.fetchRows(rows)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to fetch foreign key association: %v", err)
	}

	// 恢复原始模型
	e.model = originalModel

	return nil
}

// loadReferenceAssociation 加载引用关联
func (e *Engine) loadReferenceAssociation(mainModelValue reflect.Value, fieldValue reflect.Value, fieldType reflect.Type, foreignKey, references string, args ...any) error {
	// 引用关联通常是通过 references 指定的字段进行关联
	// 例如: belongs to 关系，使用 references 指定被引用的字段

	// 处理主模型指针
	mainModelReflectValue := mainModelValue
	if mainModelReflectValue.Kind() == reflect.Ptr {
		mainModelReflectValue = mainModelReflectValue.Elem()
	}

	mainModelType := mainModelReflectValue.Type()

	// 查找引用字段的值
	var refValue interface{}
	found := false
	for i := 0; i < mainModelType.NumField(); i++ {
		field := mainModelType.Field(i)
		// 检查字段名是否匹配，或者检查结构体标签中是否有匹配
		if strings.EqualFold(field.Name, references) || strings.EqualFold(e.getTagValue(field), references) {
			fieldVal := mainModelReflectValue.Field(i)
			refValue = fieldVal.Interface()
			found = true
			break
		}
		// 如果没有直接匹配，检查是否有类似名称的字段
		fieldSnakeCase := snakeCase(field.Name)
		if strings.EqualFold(fieldSnakeCase, references) {
			fieldVal := mainModelReflectValue.Field(i)
			refValue = fieldVal.Interface()
			found = true
			break
		}
	}

	if !found {
		// 再次遍历所有字段，尝试找到可能的引用字段
		for i := 0; i < mainModelType.NumField(); i++ {
			field := mainModelType.Field(i)
			// 检查是否有 gorm 标签包含 references
			gormTag := recursiveTags(field, types.TAG_NAME_GORM, types.TAG_NAME_SQLCA)
			if gormTag != "" {
				settings := parseTagSetting(gormTag, ";")
				// 检查是否在其他字段的 gorm 标签中定义了 references
				if fieldReferences, ok := settings["REFERENCES"]; ok {
					if strings.EqualFold(fieldReferences, references) {
						fieldVal := mainModelReflectValue.Field(i)
						refValue = fieldVal.Interface()
						found = true
						break
					}
				}
			}
		}
	}

	if !found {
		return fmt.Errorf("reference field '%s' value not found in main model", references)
	}

	// 准备关联模型
	var assocModel interface{}
	if fieldType.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			// 创建新的指针实例
			newPtr := reflect.New(fieldType.Elem())
			fieldValue.Set(newPtr)
			assocModel = newPtr.Interface()
		} else {
			assocModel = fieldValue.Interface()
		}
	} else {
		// 对于非指针类型，需要创建一个指针指向它
		assocModel = fieldValue.Addr().Interface()
	}

	// 构建查询：SELECT {fields} FROM assoc_table WHERE foreign_key_field = ref_value
	elementType := fieldType
	if fieldType.Kind() == reflect.Ptr {
		elementType = fieldType.Elem()
	}

	tableName := getTableNameFromType(elementType)

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, foreignKey)

	// 处理不同数据库的占位符
	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tableName, foreignKey)
	}

	// 临时设置模型为关联模型
	originalModel := e.model
	e.model = assocModel

	// 查询关联数据
	rows, err := e.db.Query(query, refValue)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to query reference association: %v", err)
	}
	defer rows.Close()

	// 使用fetchRows填充数据
	_, err = e.fetchRows(rows)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to fetch reference association: %v", err)
	}

	// 恢复原始模型
	e.model = originalModel

	return nil
}

// getPKValue 从结构体值中提取主键值（处理嵌套结构体如 BaseModel）
func getPKValue(modelValue reflect.Value) interface{} {
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}
	if modelValue.Kind() != reflect.Struct {
		return nil
	}
	modelType := modelValue.Type()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldVal := modelValue.Field(i)

		// 处理匿名嵌入结构体（如 BaseModel），递归查找主键
		if field.Anonymous && field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Time{}) {
			if pk := getPKValue(fieldVal); pk != nil {
				return pk
			}
		}

		// 跳过非导出字段和非主键字段
		if field.PkgPath != "" {
			continue
		}

		// 检查是否有主键标签或字段名为 ID/Id
		gormTag := recursiveTags(field, types.TAG_NAME_GORM, types.TAG_NAME_SQLCA)
		settings := parseTagSetting(gormTag, ";")
		if _, isPK := settings["PRIMARYKEY"]; isPK {
			return fieldVal.Interface()
		}
		if strings.EqualFold(field.Name, "Id") || strings.EqualFold(field.Name, "ID") {
			return fieldVal.Interface()
		}
	}
	return nil
}

// getFKFieldIndex 在结构体中查找外键字段对应的索引
func getFKFieldIndex(structType reflect.Type, foreignKey string) (int, error) {
	// 尝试直接匹配字段名（如 "UserId"）
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Time{}) {
			if idx, err := getFKFieldIndex(field.Type, foreignKey); err == nil {
				// 返回原始结构体中的第一个匹配字段
				_ = idx
			}
		}
		if strings.EqualFold(field.Name, foreignKey) {
			return i, nil
		}
	}

	// 尝试匹配蛇形命名（如 "user_id" → 查找 gorm/db 标签）
	foreignKeySnake := snakeCase(foreignKey)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Time{}) {
			if idx, err := getFKFieldIndex(field.Type, foreignKey); err == nil {
				return idx, nil
			}
		}
	}
	_ = foreignKeySnake
	return -1, fmt.Errorf("foreign key field '%s' not found in struct '%s'", foreignKey, structType.Name())
}

// loadForeignKeyAssociationBatch 批量加载外键关联（解决N+1问题）
func (e *Engine) loadForeignKeyAssociationBatch(elements []reflect.Value, pkValues []interface{}, fieldIndex int, fieldType reflect.Type, foreignKey string, args ...any) error {
	elementType := fieldType
	if fieldType.Kind() == reflect.Ptr {
		elementType = fieldType.Elem()
	}

	tableName := getTableNameFromType(elementType)
	foreignKeySnake := snakeCase(foreignKey)

	// 构建批量查询：SELECT fields FROM assoc_table WHERE fk_column IN (pk1, pk2, ...)
	placeholders := make([]string, len(pkValues))
	queryArgs := make([]interface{}, len(pkValues))
	for i, pk := range pkValues {
		placeholders[i] = "?"
		queryArgs[i] = pk
	}

	columns := getStructFields(elementType, e)
	selectFields := strings.Join(columns, ", ")
	inClause := strings.Join(placeholders, ", ")

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s IN (%s)", selectFields, tableName, foreignKeySnake, inClause)

	// 处理额外条件
	if len(args) > 0 {
		if condition, ok := args[0].(string); ok {
			query += " AND " + condition
			if len(args) > 1 {
				queryArgs = append(queryArgs, args[1:]...)
			}
		}
	}

	// 处理不同数据库的占位符
	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		inPlace := make([]string, len(pkValues))
		for i := 0; i < len(pkValues); i++ {
			inPlace[i] = "$" + strconv.Itoa(i+1)
		}
		inClause = strings.Join(inPlace, ", ")
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s IN (%s)", selectFields, tableName, foreignKeySnake, inClause)

		argOffset := len(pkValues) + 1
		if len(args) > 0 {
			if condition, ok := args[0].(string); ok {
				placeholderCount := strings.Count(condition, "?")
				query += " AND " + strings.ReplaceAll(condition, "?", fmt.Sprintf("$%d", argOffset))
				argOffset += placeholderCount
			}
		}
		queryArgs = make([]interface{}, len(pkValues))
		for i, pk := range pkValues {
			queryArgs[i] = pk
		}
		if len(args) > 0 {
			if _, ok := args[0].(string); ok {
				if len(args) > 1 {
					queryArgs = append(queryArgs, args[1:]...)
				}
			}
		}
	}

	// 创建关联模型切片
	sliceType := reflect.SliceOf(reflect.PtrTo(elementType))
	resultSlice := reflect.New(sliceType).Elem()

	// 临时设置模型为关联模型切片
	originalModel := e.model
	e.model = resultSlice.Addr().Interface()

	rows, err := e.db.Query(query, queryArgs...)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to batch query foreign key association: %v", err)
	}
	defer rows.Close()

	_, err = e.fetchRows(rows)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to batch fetch foreign key association: %v", err)
	}

	e.model = originalModel

	// 将结果按外键值映射（用于 one-to-one/has-one 关系）
	resultMap := make(map[interface{}]reflect.Value)
	for i := 0; i < resultSlice.Len(); i++ {
		result := resultSlice.Index(i)
		resultElem := result
		if resultElem.Kind() == reflect.Ptr {
			resultElem = resultElem.Elem()
		}
		fkVal := getFieldValueByName(resultElem, foreignKey)
		if fkVal != nil {
			resultMap[fkVal] = result
		}
	}

	// 分配到每个元素
	for i, element := range elements {
		pk := pkValues[i]
		fieldVal := element.Field(fieldIndex)

		if result, ok := resultMap[pk]; ok {
			if fieldVal.Kind() == reflect.Ptr {
				// 如果字段是指针，设置指向结果的指针
				fieldVal.Set(result)
			} else if fieldVal.CanSet() {
				// 如果字段是值类型
				if result.Kind() == reflect.Ptr {
					fieldVal.Set(result.Elem())
				} else {
					fieldVal.Set(result)
				}
			}
		}
	}

	return nil
}

// loadMany2ManyAssociationBatch 批量加载多对多关联（解决N+1问题）
func (e *Engine) loadMany2ManyAssociationBatch(elements []reflect.Value, pkValues []interface{}, fieldIndex int, fieldType reflect.Type, junctionTable string, args ...any) error {
	if len(elements) == 0 {
		return nil
	}

	elementType := fieldType
	if fieldType.Kind() == reflect.Ptr {
		elementType = fieldType.Elem()
	}
	if elementType.Kind() == reflect.Slice {
		elementType = elementType.Elem()
	}
	if elementType.Kind() == reflect.Ptr {
		elementType = elementType.Elem()
	}

	mainModelType := elements[0].Type()
	mainModelName := snakeCase(mainModelType.Name())
	assocModelName := snakeCase(elementType.Name())

	junctionPK1 := mainModelName + "_id"
	junctionPK2 := assocModelName + "_id"

	// 第1步：批量查询中间表，获取所有主模型与关联模型的映射关系
	placeholders := make([]string, len(pkValues))
	for i, pk := range pkValues {
		placeholders[i] = "?"
		_ = pk
	}

	// 构建中间表查询
	inClause := strings.Join(placeholders, ", ")
	junctionQuery := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s IN (%s)", junctionPK1, junctionPK2, junctionTable, junctionPK1, inClause)

	// 处理不同数据库的占位符
	junctionArgs := make([]interface{}, len(pkValues))
	copy(junctionArgs, pkValues)

	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		inPlace := make([]string, len(pkValues))
		for i := 0; i < len(pkValues); i++ {
			inPlace[i] = "$" + strconv.Itoa(i+1)
		}
		inClause = strings.Join(inPlace, ", ")
		junctionQuery = fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s IN (%s)", junctionPK1, junctionPK2, junctionTable, junctionPK1, inClause)
	}

	junctionRows, err := e.db.Query(junctionQuery, junctionArgs...)
	if err != nil {
		return fmt.Errorf("failed to batch query junction table: %v", err)
	}
	defer junctionRows.Close()

	// 收集所有关联ID并建立主模型ID到关联ID的映射
	mainToAssocIDs := make(map[uint64][]uint64)
	uniqueAssocIDs := make(map[uint64]struct{})
	var allAssocIDs []uint64

	for junctionRows.Next() {
		var mainID, assocID uint64
		if err := junctionRows.Scan(&mainID, &assocID); err != nil {
			return fmt.Errorf("failed to scan junction row: %v", err)
		}
		mainToAssocIDs[mainID] = append(mainToAssocIDs[mainID], assocID)
		if _, exists := uniqueAssocIDs[assocID]; !exists {
			uniqueAssocIDs[assocID] = struct{}{}
			allAssocIDs = append(allAssocIDs, assocID)
		}
	}

	if len(allAssocIDs) == 0 {
		// 没有关联数据，为每个元素设置空切片
		for i := range elements {
			fieldVal := elements[i].Field(fieldIndex)
			if fieldVal.Kind() == reflect.Ptr {
				newSlice := reflect.New(fieldVal.Type().Elem())
				newSlice.Elem().Set(reflect.MakeSlice(fieldVal.Type().Elem(), 0, 0))
				fieldVal.Set(newSlice)
			} else if fieldVal.Kind() == reflect.Slice {
				fieldVal.Set(reflect.MakeSlice(fieldVal.Type(), 0, 0))
			}
		}
		return nil
	}

	// 第2步：批量查询关联表
	assocPlaceholders := make([]string, len(allAssocIDs))
	for i := range allAssocIDs {
		assocPlaceholders[i] = "?"
	}
	assocInClause := strings.Join(assocPlaceholders, ", ")

	columns := getStructFields(elementType, e)
	selectFields := strings.Join(columns, ", ")
	assocTableName := getTableNameFromType(elementType)

	assocQuery := fmt.Sprintf("SELECT %s FROM %s WHERE id IN (%s)", selectFields, assocTableName, assocInClause)

	// 将 uint64 切片转换为 interface{} 切片用于 SQL 查询
	totalArgs := make([]interface{}, len(allAssocIDs))
	for i, id := range allAssocIDs {
		totalArgs[i] = id
	}

	// 处理额外条件
	if len(args) > 0 {
		if condition, ok := args[0].(string); ok {
			assocQuery += " AND " + condition
			if len(args) > 1 {
				totalArgs = append(totalArgs, args[1:]...)
			}
		}
	}

	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		inPlace := make([]string, len(allAssocIDs))
		for i := 0; i < len(allAssocIDs); i++ {
			inPlace[i] = "$" + strconv.Itoa(i+1)
		}
		assocInClause = strings.Join(inPlace, ", ")
		assocQuery = fmt.Sprintf("SELECT %s FROM %s WHERE id IN (%s)", selectFields, assocTableName, assocInClause)

		argOffset := len(allAssocIDs) + 1
		if len(args) > 0 {
			if condition, ok := args[0].(string); ok {
				placeholderCount := strings.Count(condition, "?")
				assocQuery += " AND " + strings.ReplaceAll(condition, "?", fmt.Sprintf("$%d", argOffset))
				argOffset += placeholderCount
			}
		}
		totalArgs = make([]interface{}, len(allAssocIDs))
		for i, id := range allAssocIDs {
			totalArgs[i] = id
		}
		if len(args) > 0 {
			if _, ok := args[0].(string); ok {
				if len(args) > 1 {
					totalArgs = append(totalArgs, args[1:]...)
				}
			}
		}
	}

	// 创建关联模型切片并查询
	sliceType := reflect.SliceOf(reflect.PtrTo(elementType))
	resultSlice := reflect.New(sliceType).Elem()

	originalModel := e.model
	e.model = resultSlice.Addr().Interface()

	assocRows, err := e.db.Query(assocQuery, totalArgs...)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to batch query associated data: %v", err)
	}
	defer assocRows.Close()

	_, err = e.fetchRows(assocRows)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to batch fetch associated data: %v", err)
	}

	e.model = originalModel

	// 建立关联ID到关联模型切片的映射
	assocMap := make(map[uint64]reflect.Value)
	for i := 0; i < resultSlice.Len(); i++ {
		result := resultSlice.Index(i)
		resultElem := result
		if resultElem.Kind() == reflect.Ptr {
			resultElem = resultElem.Elem()
		}
		// 查找关联模型的ID字段
		idVal := getPKValue(resultElem)
		if idVal != nil {
			// 将 ID 值转换为 uint64 类型用于 map key
			if id, ok := idVal.(uint64); ok {
				assocMap[id] = result
			}
		}
	}

	// 第3步：分配到每个元素
	for i, element := range elements {
		fieldVal := element.Field(fieldIndex)

		// 将 pkValue 转换为 uint64 类型用于 map 查找
		var pk uint64
		if pkv, ok := pkValues[i].(uint64); ok {
			pk = pkv
		}

		// 获取该主模型关联的所有关联ID
		assocIDs, hasAssoc := mainToAssocIDs[pk]
		if !hasAssoc || len(assocIDs) == 0 {
			// 没有关联数据，设置空切片
			if fieldVal.Kind() == reflect.Ptr {
				newSlice := reflect.New(fieldVal.Type().Elem())
				newSlice.Elem().Set(reflect.MakeSlice(fieldVal.Type().Elem(), 0, 0))
				fieldVal.Set(newSlice)
			} else if fieldVal.Kind() == reflect.Slice {
				fieldVal.Set(reflect.MakeSlice(fieldVal.Type(), 0, 0))
			}
			continue
		}

		// 创建结果切片
		elemSliceType := fieldVal.Type()
		if elemSliceType.Kind() == reflect.Ptr {
			elemSliceType = elemSliceType.Elem()
		}

		elemSlice := reflect.MakeSlice(elemSliceType, 0, len(assocIDs))
		for _, assocID := range assocIDs {
			if result, ok := assocMap[assocID]; ok {
				if elemSliceType.Elem().Kind() == reflect.Ptr {
					elemSlice = reflect.Append(elemSlice, result)
				} else {
					if result.Kind() == reflect.Ptr {
						elemSlice = reflect.Append(elemSlice, result.Elem())
					} else {
						elemSlice = reflect.Append(elemSlice, result)
					}
				}
			}
		}

		if fieldVal.Kind() == reflect.Ptr {
			newSlice := reflect.New(fieldVal.Type().Elem())
			newSlice.Elem().Set(elemSlice)
			fieldVal.Set(newSlice)
		} else {
			fieldVal.Set(elemSlice)
		}
	}

	return nil
}

// loadReferenceAssociationBatch 批量加载引用关联（解决N+1问题）
func (e *Engine) loadReferenceAssociationBatch(elements []reflect.Value, pkValues []interface{}, fieldIndex int, fieldType reflect.Type, foreignKey, references string, args ...any) error {
	elementType := fieldType
	if fieldType.Kind() == reflect.Ptr {
		elementType = fieldType.Elem()
	}

	// 收集所有元素的引用字段值
	refValues := make([]interface{}, len(elements))
	mainModelType := elements[0].Type()
	refFieldIdx := -1

	for i := 0; i < mainModelType.NumField(); i++ {
		field := mainModelType.Field(i)
		if strings.EqualFold(field.Name, references) {
			refFieldIdx = i
			break
		}
		fieldSnakeCase := snakeCase(field.Name)
		if strings.EqualFold(fieldSnakeCase, references) {
			refFieldIdx = i
			break
		}
	}

	if refFieldIdx == -1 {
		// 没有找到引用字段，降级为逐个加载
		for i, element := range elements {
			if err := e.loadReferenceAssociation(element, elements[i].Field(fieldIndex), fieldType, foreignKey, references, args...); err != nil {
				return err
			}
			_ = i
		}
		return nil
	}

	for i, element := range elements {
		refValues[i] = element.Field(refFieldIdx).Interface()
	}

	tableName := getTableNameFromType(elementType)
	foreignKeySnake := snakeCase(foreignKey)

	// 构建批量查询
	placeholders := make([]string, len(refValues))
	queryArgs := make([]interface{}, len(refValues))
	for i, rv := range refValues {
		placeholders[i] = "?"
		queryArgs[i] = rv
	}

	columns := getStructFields(elementType, e)
	selectFields := strings.Join(columns, ", ")
	inClause := strings.Join(placeholders, ", ")

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s IN (%s)", selectFields, tableName, foreignKeySnake, inClause)

	if len(args) > 0 {
		if condition, ok := args[0].(string); ok {
			query += " AND " + condition
			if len(args) > 1 {
				queryArgs = append(queryArgs, args[1:]...)
			}
		}
	}

	switch e.adapterType {
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		inPlace := make([]string, len(refValues))
		for i := 0; i < len(refValues); i++ {
			inPlace[i] = "$" + strconv.Itoa(i+1)
		}
		inClause = strings.Join(inPlace, ", ")
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s IN (%s)", selectFields, tableName, foreignKeySnake, inClause)

		argOffset := len(refValues) + 1
		if len(args) > 0 {
			if condition, ok := args[0].(string); ok {
				placeholderCount := strings.Count(condition, "?")
				query += " AND " + strings.ReplaceAll(condition, "?", fmt.Sprintf("$%d", argOffset))
				argOffset += placeholderCount
			}
		}
		queryArgs = make([]interface{}, len(refValues))
		for i, rv := range refValues {
			queryArgs[i] = rv
		}
		if len(args) > 0 {
			if _, ok := args[0].(string); ok {
				if len(args) > 1 {
					queryArgs = append(queryArgs, args[1:]...)
				}
			}
		}
	}

	// 创建关联模型切片
	sliceType := reflect.SliceOf(reflect.PtrTo(elementType))
	resultSlice := reflect.New(sliceType).Elem()

	originalModel := e.model
	e.model = resultSlice.Addr().Interface()

	rows, err := e.db.Query(query, queryArgs...)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to batch query reference association: %v", err)
	}
	defer rows.Close()

	_, err = e.fetchRows(rows)
	if err != nil {
		e.model = originalModel
		return fmt.Errorf("failed to batch fetch reference association: %v", err)
	}

	e.model = originalModel

	// 将结果按外键值映射
	resultMap := make(map[interface{}]reflect.Value)
	for i := 0; i < resultSlice.Len(); i++ {
		result := resultSlice.Index(i)
		resultElem := result
		if resultElem.Kind() == reflect.Ptr {
			resultElem = resultElem.Elem()
		}
		fkVal := getFieldValueByName(resultElem, foreignKey)
		if fkVal != nil {
			resultMap[fkVal] = result
		}
	}

	// 分配到每个元素
	for i, element := range elements {
		refVal := refValues[i]
		fieldVal := element.Field(fieldIndex)

		if result, ok := resultMap[refVal]; ok {
			if fieldVal.Kind() == reflect.Ptr {
				fieldVal.Set(result)
			} else if fieldVal.CanSet() {
				if result.Kind() == reflect.Ptr {
					fieldVal.Set(result.Elem())
				} else {
					fieldVal.Set(result)
				}
			}
		}
	}

	return nil
}

// getFieldValueByName 根据 Go 字段名（或蛇形命名的 db 标签）从结构体中查找字段值
func getFieldValueByName(modelValue reflect.Value, fieldName string) interface{} {
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}
	if modelValue.Kind() != reflect.Struct {
		return nil
	}
	modelType := modelValue.Type()

	// 直接匹配字段名
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if strings.EqualFold(field.Name, fieldName) {
			return modelValue.Field(i).Interface()
		}
	}

	// 匹配蛇形命名对应的 db 标签
	fieldNameSnake := snakeCase(fieldName)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		gormTag := recursiveTags(field, types.TAG_NAME_GORM, types.TAG_NAME_SQLCA)
		settings := parseTagSetting(gormTag, ";")
		if colName, ok := settings["COLUMN"]; ok && strings.EqualFold(colName, fieldNameSnake) {
			return modelValue.Field(i).Interface()
		}
	}

	// 处理匿名嵌入结构体（如 BaseModel）中寻找字段
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Time{}) {
			if val := getFieldValueByName(modelValue.Field(i), fieldName); val != nil {
				return val
			}
		}
	}

	// 最后尝试直接匹配索引（非导出字段等）
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Time{}) {
			_ = field
		}
	}

	return nil
}

// getTableNameFromType 从类型获取表名
func getTableNameFromType(t reflect.Type) string {
	// 尝试调用TableName方法
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 创建一个临时实例来调用TableName方法
	tempInstance := reflect.New(t).Interface()
	if tabler, ok := tempInstance.(interface{ TableName() string }); ok {
		return tabler.TableName()
	}

	// 如果没有TableName方法，使用蛇形命名规则
	return snakeCase(t.Name())
}

// snakeCase 转换为蛇形命名
func snakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return strings.ToLower(string(result))
}

// convertMapInterfaceToMapString 将 map[string]interface{} 转换为 map[string]string
func convertMapInterfaceToMapString(input map[string]interface{}) map[string]string {
	result := make(map[string]string, len(input))
	for k, v := range input {
		if v == nil {
			result[k] = ""
		} else {
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	return result
}

// getStructFields 获取结构体的所有数据库字段名
func getStructFields(structType reflect.Type, e *Engine) []string {
	var columns []string

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// 跳过非导出字段
		if field.PkgPath != "" {
			continue
		}

		// 检查是否是嵌套结构体
		if field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Time{}) {
			// 递归处理嵌套结构体的字段
			nestedColumns := getStructFields(field.Type, e)
			columns = append(columns, nestedColumns...)
			continue
		}

		// 获取字段的数据库标签
		tagName := e.getTagValue(field)
		if tagName == "" {
			// 如果没有db标签，使用字段名的蛇形格式
			tagName = snakeCase(field.Name)
		}

		// 如果标签值为"-"表示忽略此字段
		if tagName == "-" {
			continue
		}

		columns = append(columns, tagName)
	}

	return columns
}

// isInOperator 检查字符串是否为常见SQL操作符
func isInOperator(op string) bool {
	switch strings.ToUpper(op) {
	case "=", "!=", "<>", "<", "<=", ">", ">=", "LIKE", "NOT LIKE", "IN", "NOT IN", "EXISTS":
		return true
	default:
		return false
	}
}
