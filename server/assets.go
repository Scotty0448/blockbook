package server

import (
	"strconv"
	"encoding/binary"
 	"github.com/martinboehm/btcutil/base58"
)

type AssetData struct {
  Header     string
  TxType     string
  Name       string
  Quantity   string
  Units      uint8
  Reissuable bool
  HashType   string
  Hash       string
}

func IsAssetScript(pkScript []byte) (bool) {
  return len(pkScript) > 25 && pkScript[25] == 0xc0
}

func ExtractAssetData(pkScript []byte) (AssetData) {
  var data AssetData

  script_length := pkScript[26]
  name_length   := pkScript[31]

  data.Header = string(pkScript[27:30])

  switch pkScript[30] {
    case 0x6f:
      data.TxType = "New Ownership"
    case 0x71:
      data.TxType = "New"
    case 0x72:
      data.TxType = "Reissue"
    case 0x74:
      data.TxType = "Transfer"
  }

  data.Name = string(pkScript[32:32+name_length])

  if 27+script_length > 32+name_length {
    quantity := binary.LittleEndian.Uint64(pkScript[32+name_length:40+name_length])
    data.Quantity = strconv.FormatFloat(float64(quantity/100000000), 'f', -1, 64)
  } else {
    data.Quantity = "1"
  }

  if 27+script_length > 40+name_length {
    data.Units = uint8(pkScript[40+name_length])
  }

  if 27+script_length > 41+name_length {
    data.Reissuable = pkScript[41+name_length] == 1
    switch pkScript[42+name_length] {
      case 0x00:
        data.HashType = ""
        data.Hash     = ""
      case 0x01:
        data.HashType = "IPFS"
        data.Hash     = base58.Encode(pkScript[43+name_length:26+script_length])
    }
  }

  return data
}
