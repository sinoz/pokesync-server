package game

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// AssetConfig holds configurations specific to game assets.
type AssetConfig struct {
	ItemDirectory    string
	NpcDirectory     string
	MonsterDirectory string
	ObjectDirectory  string
	WorldDirectory   string
}

// ItemDescriptor describes an item.
type ItemDescriptor struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	NitroLabelID int    `json:"nitroLabelId"`
	StorePrice   int    `json:"storePrice"`
}

// ItemConfig contains a collection of item config entries.
type ItemConfig struct {
	Descriptors []ItemDescriptor
}

// NpcDescriptor describes a npc.
type NpcDescriptor struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	NitroLabelID int    `json:"nitroLabelId"`
}

// NpcConfig contains a collection of npc config entries.
type NpcConfig struct {
	Descriptors []NpcDescriptor
}

// TypePair is the pair of monster types.
type TypePair struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

// MonsterDescriptor describes a monster.
type MonsterDescriptor struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	NitroLabelID  int       `json:"nitroLabelId"`
	Height        int       `json:"height"`
	Weight        int       `json:"weight"`
	Types         TypePair  `json:"types"`
	BaseExp       int       `json:"baseExperience"`
	BaseStats     [7]int    `json:"baseStats"`
	Abilities     [8]string `json:"abilities"`
	PreviousEvo   int       `json:"previousEvolution"`
	GenderRate    int       `json:"genderRate"`
	CaptureRate   int       `json:"captureRate"`
	BaseHappiness int       `json:"baseHappiness"`
}

// MonsterConfig contains a collection of monster descriptors.
type MonsterConfig struct {
	Descriptors []MonsterDescriptor
}

// ObjectDescriptor describes an object.
type ObjectDescriptor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ObjectConfig contains a collection of object config entries.
type ObjectConfig struct {
	Descriptor []ObjectDescriptor
}

// MapIndex is a single index entry of a map on the global grid.
type MapIndex struct {
	MapX    int `json:"mapX"`
	MapZ    int `json:"mapZ"`
	RenderX int `json:"renderX"`
	RenderZ int `json:"renderZ"`
}

// RegionIndex is a single entry of a region on the global grid.
type RegionIndex struct {
	Label string     `json:"label"`
	Maps  []MapIndex `json:"maps"`
}

// WorldConfig contains a collection of RegionIndex entries.
type WorldConfig struct {
	Width   int           `json:"width"`
	Length  int           `json:"length"`
	Regions []RegionIndex `json:"regions"`
}

// AssetBundle hold all of the assets the game will need.
type AssetBundle struct {
	Items    *ItemConfig
	Npcs     *NpcConfig
	Monsters *MonsterConfig
	Objects  *ObjectConfig
	World    *WorldConfig
	Grid     *Grid
}

// LoadAssetBundle loads all of the assets that the game requires.
func LoadAssetBundle(config AssetConfig) (*AssetBundle, error) {
	itemConfig, err := LoadItemConfigsAt(config.ItemDirectory)
	if err != nil {
		return nil, err
	}

	npcConfig, err := LoadNpcConfigsAt(config.NpcDirectory)
	if err != nil {
		return nil, err
	}

	monsterConfig, err := LoadMonsterConfigsAt(config.MonsterDirectory)
	if err != nil {
		return nil, err
	}

	worldConfig, err := LoadWorldConfigAt(config.WorldDirectory + "/world.json")
	if err != nil {
		return nil, err
	}

	grid, err := LoadGridFromConfig(*worldConfig)
	if err != nil {
		return nil, err
	}

	return &AssetBundle{
		Items:    itemConfig,
		Npcs:     npcConfig,
		Monsters: monsterConfig,
		World:    worldConfig,
		Grid:     grid,
	}, nil
}

// LoadItemConfigsAt loads all of the item config files found within the
// specified directory. May return an error.
func LoadItemConfigsAt(directory string) (*ItemConfig, error) {
	fileInfo, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	itemCount := len(fileInfo)

	config := new(ItemConfig)
	config.Descriptors = make([]ItemDescriptor, itemCount)

	for i := 0; i < itemCount; i++ {
		fileBytes, err := ioutil.ReadFile(directory + "/" + strconv.Itoa(i) + ".json")
		if err != nil {
			return nil, err
		}

		descriptor := &ItemDescriptor{}
		if err := json.Unmarshal(fileBytes, descriptor); err != nil {
			return nil, err
		}

		config.Descriptors[i] = *descriptor
	}

	return config, nil
}

// LoadMonsterConfigsAt loads the given amount of monster config files from
// the specified directory. May return an error.
func LoadMonsterConfigsAt(directory string) (*MonsterConfig, error) {
	fileInfo, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	monsterCount := len(fileInfo)

	config := new(MonsterConfig)
	config.Descriptors = make([]MonsterDescriptor, monsterCount)

	for i := 0; i < monsterCount; i++ {
		fileBytes, err := ioutil.ReadFile(directory + "/" + strconv.Itoa(i) + ".json")
		if err != nil {
			return nil, err
		}

		descriptor := &MonsterDescriptor{}
		if err := json.Unmarshal(fileBytes, descriptor); err != nil {
			return nil, err
		}

		config.Descriptors[i] = *descriptor
	}

	return config, nil
}

// Count returns the amount of loaded descriptors this MonsterConfig holds.
func (config *MonsterConfig) Count() int {
	return len(config.Descriptors)
}

// LoadNpcConfigsAt loads the given amount of item config files from
// the specified directory. May return an error.
func LoadNpcConfigsAt(directory string) (*NpcConfig, error) {
	fileInfo, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	npcCount := len(fileInfo)

	config := new(NpcConfig)
	config.Descriptors = make([]NpcDescriptor, npcCount)

	for i := 0; i < npcCount; i++ {
		fileBytes, err := ioutil.ReadFile(directory + "/" + strconv.Itoa(i) + ".json")
		if err != nil {
			return nil, err
		}

		descriptor := &NpcDescriptor{}
		if err := json.Unmarshal(fileBytes, descriptor); err != nil {
			return nil, err
		}

		config.Descriptors[i] = *descriptor
	}

	return config, nil
}

// LoadWorldConfigAt loads a WorldConfig from a world config file
// at the specified path. May return an error.
func LoadWorldConfigAt(path string) (*WorldConfig, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &WorldConfig{}
	if err := json.Unmarshal(fileBytes, config); err != nil {
		return nil, err
	}

	return config, nil
}

// Count returns the amount of loaded item descriptors this ItemConfig holds.
func (config *ItemConfig) Count() int {
	if config.Descriptors == nil {
		return 0
	}

	return len(config.Descriptors)
}

// Count returns the amount of loaded npc descriptors this NpcConfig holds.
func (config *NpcConfig) Count() int {
	if config.Descriptors == nil {
		return 0
	}

	return len(config.Descriptors)
}
