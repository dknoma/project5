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

func DecodeEquipmentFromJson(jsonString string) (Equipment, error) {
	jsonBytes := []byte(jsonString)
	// Unmarshal the json bytes into a new key:value map
	var equipmentMap map[string]interface{}
	err := json.Unmarshal(jsonBytes, &equipmentMap)
	if err != nil {
		fmt.Println(err.Error())
		return Equipment{}, err
	}
	// Create new block to insert unmarshalled values into
	var e Equipment
	e.Name = equipmentMap["name"].(string)
	e.Id = int32(equipmentMap["id"].(float64))
	e.Owner = int32(equipmentMap["owner"].(float64))
	e.Description = equipmentMap["description"].(string)

	// convert stats string into bytes into struct
	statBytes := []byte(equipmentMap["stats"].(string))
	var statMap map[string]interface{}
	err = json.Unmarshal(statBytes, &statMap)
	if err != nil {
		fmt.Println(err.Error())
		return Equipment{}, err
	}
	var es EquipmentStats
	es.Level = int32(statMap["level"].(float64))
	es.Atk = int32(statMap["atk"].(float64))
	es.Def = int32(statMap["def"].(float64))
	// Insert stats into equipement
	e.Stats = es
	return e, nil
}
