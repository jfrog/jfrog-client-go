package utils

import (
	"bytes"
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"io"
)

func EncodeBomToJson(bom *cyclonedx.BOM) ([]byte, error) {
	// Encode the BOM to JSON format
	var buf bytes.Buffer
	var writer io.Writer = &buf
	encoder := cyclonedx.NewBOMEncoder(writer, cyclonedx.BOMFileFormatJSON)
	if err := encoder.Encode(bom); err != nil {
		return nil, errorutils.CheckErrorf("failed to encode CycloneDX BOM: %s", err.Error())
	}
	return buf.Bytes(), nil
}

func DecodeBomFromJson(bomJson []byte) (*cyclonedx.BOM, error) {
	// Decode the BOM back to a CycloneDX BOM object
	reader := bytes.NewReader(bomJson)
	decoder := cyclonedx.NewBOMDecoder(reader, cyclonedx.BOMFileFormatJSON)
	bom := &cyclonedx.BOM{}
	if err := decoder.Decode(bom); err != nil {
		return nil, errorutils.CheckErrorf("failed to decode CycloneDX BOM: %s", err.Error())
	}
	return bom, nil
}
