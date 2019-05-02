package gamedata

import (
	"encoding/json"
	"fmt"
)

type EquipmentStats struct {
	Level int32 `json:"level"`
	Atk   int32 `json:"atk"`
	Def   int32 `json:"def"`
}

type Equipment struct {
	Name        string         `json:"name"`        // Name of the equipment
	Id          int32          `json:"id"`          // Id of the equipment
	Owner       int32          `json:"owner"`       // Id of the equipments owner
	Description string         `json:"description"` // Description of the equipment
	Stats       EquipmentStats `json:"stats"`       // Equipment stats
}

func New(name string, id int32, owner int32, description string, level, atk, def int32) Equipment {
	return Equipment{name, id, owner, description, EquipmentStats{level, atk, def}}
}

func (e *Equipment) EncodeEquipmentToJson() (string, error) {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	mptJson := string(jsonBytes)
	mptJson = mptJson[1 : len(mptJson)-1]
	jsonOut := fmt.Sprintf("{\"name\": \"%v\",\"id\": %v,\"owner\": %v,\"description\": \"%v\",\"stats\": {\"level\": %v,\"atk\": %v,\"def\": %v}}",
		e.Name, e.Id, e.Owner, e.Description, e.Stats.Level, e.Stats.Atk, e.Stats.Def)
	isValid := json.Valid([]byte(jsonOut))
	if !isValid {
		fmt.Println(err.Error())
		return "", err
	}
	return jsonOut, nil
}
