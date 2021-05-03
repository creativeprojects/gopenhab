package api

import (
	"log"
	"strings"
)

type ItemType string

const (
	itemTypeNumberWithSubtype = "Number:"
)

const (
	ItemTypeUnknown       ItemType = "Unknown"
	ItemTypeContact       ItemType = "Contact"
	ItemTypeColor         ItemType = "Color"
	ItemTypeDateTime      ItemType = "DateTime"
	ItemTypeDimmer        ItemType = "Dimmer"
	ItemTypeGroup         ItemType = "Group"
	ItemTypeImage         ItemType = "Image"
	ItemTypeLocation      ItemType = "Location"
	ItemTypeNumber        ItemType = "Number"
	ItemTypeRollershutter ItemType = "Rollershutter"
	ItemTypeString        ItemType = "String"
	ItemTypeSwitch        ItemType = "Switch"
)

func GetItemType(s string) (ItemType, string) {
	// Special case of Number that can have a subtype (which is not kept anywhere, or at least for now)
	if strings.HasPrefix(s, itemTypeNumberWithSubtype) {
		return ItemTypeNumber, strings.TrimPrefix(s, itemTypeNumberWithSubtype)
	}

	switch s {
	default:
		log.Printf("unknown type: %s", s)
		return ItemTypeUnknown, s
	case string(ItemTypeColor):
		return ItemTypeColor, ""
	case string(ItemTypeContact):
		return ItemTypeContact, ""
	case string(ItemTypeDateTime):
		return ItemTypeDateTime, ""
	case string(ItemTypeDimmer):
		return ItemTypeDimmer, ""
	case string(ItemTypeGroup):
		return ItemTypeGroup, ""
	case string(ItemTypeImage):
		return ItemTypeImage, ""
	case string(ItemTypeLocation):
		return ItemTypeLocation, ""
	case string(ItemTypeNumber):
		return ItemTypeNumber, ""
	case string(ItemTypeRollershutter):
		return ItemTypeRollershutter, ""
	case string(ItemTypeString):
		return ItemTypeString, ""
	case string(ItemTypeSwitch):
		return ItemTypeSwitch, ""
	}
}
