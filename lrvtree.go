package sqlca

// LeftRightValTree 左右值树
type LeftRightValTree struct {
	LeftVal  int                  `gorm:"column:left_val;index"`                   // 左值
	RightVal int                  `gorm:"column:right_val;index"`                  // 右值
	Depth    int                  `gorm:"column:depth;"`                           // 深度
	ParentId int64                `gorm:"column:parent_id;default:0"`              // 父ID
	RootPath KindBaseSlice[int64] `gorm:"column:root_path;type:json;default:null"` // 路径
}

// 创建左右值树根节点
func NewLrvRoot() LeftRightValTree {
	return LeftRightValTree{
		LeftVal:  1,
		RightVal: 2,
		Depth:    0,
	}
}

// 创建左右值树子节点
func (lr *LeftRightValTree) NewChild(parentId int64) LeftRightValTree {
	nlr := &LeftRightValTree{}
	if parentId == 0 { // 写入节点为根节点
		lr.LeftVal, lr.RightVal, lr.Depth = 1, 2, 0
		return *lr
	}

	// 根据父节点的左右值确定当前节点的左右值
	nlr.LeftVal = lr.GetRightVal()
	nlr.RightVal = nlr.GetLeftVal() + 1
	nlr.Depth = lr.GetDepth() + 1
	nlr.ParentId = parentId
	nlr.RootPath = append(lr.GetRootPath(), parentId)
	return *nlr
}

func (lr *LeftRightValTree) GetLeftVal() int {
	return lr.LeftVal
}
func (lr *LeftRightValTree) GetRightVal() int {
	return lr.RightVal
}
func (lr *LeftRightValTree) GetDepth() int {
	return lr.Depth
}
func (lr *LeftRightValTree) GetParentId() int64 {
	return lr.ParentId
}

func (lr *LeftRightValTree) GetRootPath() []int64 {
	return lr.RootPath
}
