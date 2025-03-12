package sqlca

import (
	"github.com/civet148/log"
	"reflect"
)

type BeforeCreateInterface interface {
	BeforeCreate(db *Engine) error
}

type AfterCreateInterface interface {
	AfterCreate(db *Engine) error
}

type BeforeUpdateInterface interface {
	BeforeUpdate(db *Engine) error
}

type AfterUpdateInterface interface {
	AfterUpdate(db *Engine) error
}

type BeforeDeleteInterface interface {
	BeforeDelete(db *Engine) error
}

type AfterDeleteInterface interface {
	AfterDelete(db *Engine) error
}

type hookMethods struct {
	beforeCreates []BeforeCreateInterface // before create
	beforeUpdates []BeforeUpdateInterface // before update
	beforeDeletes []BeforeDeleteInterface // before delete
	afterCreates  []AfterCreateInterface  // after create
	afterUpdates  []AfterUpdateInterface  // after update
	afterDeletes  []AfterDeleteInterface  // after delete
}

// 获取Hook方法
func getHookMethods(obj interface{}) *hookMethods {
	var methods = &hookMethods{}
	if obj == nil {
		return methods
	}
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	//log.Infof("obj type %v", typ)
	for {
		if typ.Kind() != reflect.Ptr { // pointer type
			break
		}
		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Struct:
		{
			parseStructHookMethods(methods, typ, val)
		}
	case reflect.Slice:
		{
			for i := 0; i < val.Len(); i++ {
				elemVal := val.Index(i)
				elemTyp := elemVal.Type()
				if elemTyp.Kind() == reflect.Ptr {
					elemTyp = elemTyp.Elem()
					elemVal = elemVal.Elem()
				}
				if elemTyp.Kind() == reflect.Struct {
					parseStructHookMethods(methods, elemTyp, elemVal)
				}
			}
		}
	default:
	}
	return methods
}

// parse struct fields
func parseStructHookMethods(methods *hookMethods, typ reflect.Type, val reflect.Value) {
	if !val.CanAddr() {
		//log.Errorf("type %v can't get address", typ)
		return
	}
	if typ.Kind() == reflect.Struct {
		val = val.Addr()
	}
	if hook, ok := val.Interface().(BeforeCreateInterface); ok {
		methods.beforeCreates = append(methods.beforeCreates, hook)
	}
	if hook, ok := val.Interface().(BeforeUpdateInterface); ok {
		methods.beforeUpdates = append(methods.beforeUpdates, hook)
	}
	if hook, ok := val.Interface().(BeforeDeleteInterface); ok {
		methods.beforeDeletes = append(methods.beforeDeletes, hook)
	}
	if hook, ok := val.Interface().(AfterCreateInterface); ok {
		methods.afterCreates = append(methods.afterCreates, hook)
	}
	if hook, ok := val.Interface().(AfterUpdateInterface); ok {
		methods.afterUpdates = append(methods.afterUpdates, hook)
	}
	if hook, ok := val.Interface().(AfterDeleteInterface); ok {
		methods.afterDeletes = append(methods.afterDeletes, hook)
	}
}

func (e *Engine) cleanHooks() {
	e.hookMethods = nil
}

func (e *Engine) setHooks() *Engine {
	e.hookMethods = getHookMethods(e.model)
	return e
}

func (e *Engine) execBeforeCreateHooks() (err error) {
	if e.hookMethods != nil {
		for _, hook := range e.hookMethods.beforeCreates {
			if hook == nil {
				continue
			}
			if err = hook.BeforeCreate(e.clone()); err != nil {
				return log.Errorf(err.Error())
			}
		}
	}
	return nil
}

func (e *Engine) execBeforeUpdateHooks() (err error) {
	if e.hookMethods != nil {
		for _, hook := range e.hookMethods.beforeUpdates {
			if hook == nil {
				continue
			}
			if err = hook.BeforeUpdate(e.clone()); err != nil {
				return log.Errorf(err.Error())
			}
		}
	}
	return nil
}

func (e *Engine) execBeforeDeleteHooks() (err error) {
	if e.hookMethods != nil {
		for _, hook := range e.hookMethods.beforeDeletes {
			if hook == nil {
				continue
			}
			if err = hook.BeforeDelete(e.clone()); err != nil {
				return log.Errorf(err.Error())
			}
		}
	}
	return nil
}

func (e *Engine) execAfterCreateHooks() (err error) {
	if e.hookMethods != nil {
		for _, hook := range e.hookMethods.afterCreates {
			if hook == nil {
				continue
			}
			if err = hook.AfterCreate(e.clone()); err != nil {
				return log.Errorf(err.Error())
			}
		}
	}
	return nil
}

func (e *Engine) execAfterUpdateHooks() (err error) {
	if e.hookMethods != nil {
		for _, hook := range e.hookMethods.afterUpdates {
			if hook == nil {
				continue
			}
			if err = hook.AfterUpdate(e.clone()); err != nil {
				return log.Errorf(err.Error())
			}
		}
	}
	return nil
}

func (e *Engine) execAfterDeleteHooks() (err error) {
	if e.hookMethods != nil {
		for _, hook := range e.hookMethods.afterDeletes {
			if hook == nil {
				continue
			}
			if err = hook.AfterDelete(e.clone()); err != nil {
				return log.Errorf(err.Error())
			}
		}
	}
	return nil
}

