package nchi

import (
	"encoding/json"
	"encoding/xml"

	"github.com/muir/nvelope"

	"github.com/julienschmidt/httprouter"
)

// DecodeJSON is is a pre-defined special nject.Provider created with
// nvelope.GenerateDecoder for decoding JSON requests.  Use it with the
// other features of https://github.com/muir/nvelope . DecodeJSON must
// be paired with nvelope.ReadBody to actually decode JSON.
var DecodeJSON = nvelope.GenerateDecoder(
	nvelope.WithDecoder("application/json", json.Unmarshal),
	nvelope.WithDefaultContentType("application/json"),
	nvelope.WithPathVarsFunction(func(p httprouter.Params) nvelope.RouteVarLookup {
		return p.ByName
	}),
)

// DecodeXML is is a pre-defined special nject.Provider created with
// nvelope.GenerateDecoder for decoding XML requests.Use it with the
// other features of https://github.com/muir/nvelope .  DecodeXML must be
// paired with nvelope.ReadBody to actually decode XML.
var DecodeXML = nvelope.GenerateDecoder(
	nvelope.WithDecoder("application/xml", xml.Unmarshal),
	nvelope.WithDefaultContentType("application/xml"),
	nvelope.WithPathVarsFunction(func(p httprouter.Params) nvelope.RouteVarLookup {
		return p.ByName
	}),
)
