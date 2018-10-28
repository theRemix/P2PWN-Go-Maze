package main

import (
	"encoding/hex"
	"flag"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func loadTTF(inlinedFont string, size float64) (font.Face, error) {
	bytes, err := hex.DecodeString(inlinedFont)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func setConfig(configPtr *string, flagName string, defaultVal string, help string) {
	flag.StringVar(configPtr, flagName, defaultVal, help)

	if val, ok := os.LookupEnv(flagName); ok {
		*configPtr = val
	}
}
