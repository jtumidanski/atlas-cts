package portal

type Model struct {
	id          uint32
	name        string
	target      string
	theType     uint8
	x           int16
	y           int16
	targetMapId uint32
	scriptName  string
}

func (p Model) ScriptName() string {
	return p.scriptName
}

func (p Model) TargetMapId() uint32 {
	return p.targetMapId
}

func (p Model) Target() string {
	return p.target
}

func (p Model) Id() uint32 {
	return p.id
}

func (p Model) Name() string {
	return p.name
}

func NewPortalModel(id uint32, name string, target string, targetMapId uint32, theType uint8, x int16, y int16, scriptName string) Model {
	return Model{
		id:          id,
		name:        name,
		target:      target,
		targetMapId: targetMapId,
		theType:     theType,
		x:           x,
		y:           y,
		scriptName:  scriptName,
	}
}
