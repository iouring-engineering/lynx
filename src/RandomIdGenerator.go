package main

import (
	"errors"
	"fmt"
	"hash/fnv"
	"net"
	"os"
	"strings"
	"time"
)

var EPOCH_BITS int = 21
var NODE_ID_BITS int = 5
var SEQUENCE_BITS int = 5
var maxNodeId int64 = (1 << NODE_ID_BITS) - 1
var maxSequence int64 = (1 << SEQUENCE_BITS) - 1
var defaultCustomEpoch int64 = 1420070400000
var lastTimestamp int64 = -1
var sequence int64 = 0

type RandomIdGenerator interface {
	IdGenerator()
	IdGeneratorNode(nodeId int64)
	IdGeneratorNodeCustom(nodeId int64, customEpoch int64) error
	GetId() (int64, error)
	SetNodeBitsLength(value int)
	SetSequenceBitsLength(value int)
	SetEpochBitsLength(value int)
	SetCustomEpoch(value int64)
	ParseId(id int64) []int64
}

type RandomID struct {
	nodeId      int64
	customEpoch int64
}

func (random *RandomID) SetNodeBitsLength(value int) {
	NODE_ID_BITS = value
	maxNodeId = (1 << NODE_ID_BITS) - 1
}

func (random *RandomID) SetSequenceBitsLength(value int) {
	SEQUENCE_BITS = value
	maxNodeId = (1 << NODE_ID_BITS) - 1
}

func (random *RandomID) SetEpochBitsLength(value int) {
	EPOCH_BITS = value
}

func (random *RandomID) SetCustomEpoch(value int64) {
	random.customEpoch = value
}

// Create Snowflake with a nodeId
func (random *RandomID) IdGeneratorNode(nodeId int64) {
	random.IdGeneratorNodeCustom(nodeId, defaultCustomEpoch)
}

// Create Snowflake with a nodeId and custom epoch
func (random *RandomID) IdGeneratorNodeCustom(nodeId int64, customEpoch int64) error {
	if nodeId < 0 || nodeId > maxNodeId {
		errMsg := fmt.Sprintf("NodeId must be between %d and %d", 0, maxNodeId)
		return errors.New(errMsg)
	}
	random.nodeId = nodeId
	random.customEpoch = customEpoch
	return nil
}

// Let Snowflake generate a nodeId
func (random *RandomID) IdGenerator() {
	random.nodeId = random.generateNodeId()
	random.customEpoch = defaultCustomEpoch
}

func (random *RandomID) generateNodeId() int64 {
	var sb strings.Builder
	interfaces, err := net.Interfaces()
	if err != nil {
		Logger.Info(err.Error())
	}
	for _, i := range interfaces {
		var addrs []byte = i.HardwareAddr
		for _, mac := range addrs {
			sb.WriteString(fmt.Sprintf("%x", mac))
		}
	}
	sb.WriteString(fmt.Sprintf("%d", os.Getpid()))
	nodeId := hashCode(sb.String())
	nodeId = nodeId & maxNodeId
	return nodeId
}

func hashCode(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

func (random *RandomID) GetId() (int64, error) {
	var currentTimestamp int64 = random.getTimestamp()
	if currentTimestamp < lastTimestamp {
		return 0, fmt.Errorf("invalid system clock time")
	}
	if currentTimestamp == lastTimestamp {
		sequence = (sequence + 1) & maxSequence
		if sequence == 0 {
			currentTimestamp = random.waitNextMillis(currentTimestamp)
		}
	} else {
		// reset sequence to start with zero for the next millisecond
		sequence = 0
	}
	lastTimestamp = currentTimestamp

	randomId := (currentTimestamp<<(NODE_ID_BITS+SEQUENCE_BITS) |
		(random.nodeId << int64(SEQUENCE_BITS)) | sequence)
	return randomId, nil
}

func (random *RandomID) waitNextMillis(currentTimestamp int64) int64 {
	for currentTimestamp == lastTimestamp {
		currentTimestamp = random.getTimestamp()
	}
	return currentTimestamp
}

// Get current timestamp in milliseconds, adjust for the custom epoch.
func (random *RandomID) getTimestamp() int64 {
	tUnixMilli := time.Now().UnixNano() / int64(time.Millisecond)
	currentTime := tUnixMilli - defaultCustomEpoch
	return currentTime
}

func (random *RandomID) ParseId(id int64) []int64 {
	maskNodeId := ((1 << NODE_ID_BITS) - 1) << SEQUENCE_BITS
	maskSequence := (1 << SEQUENCE_BITS) - 1

	timestamp := (id >> (NODE_ID_BITS + SEQUENCE_BITS)) + random.customEpoch
	nodeId := (id & int64(maskNodeId)) >> SEQUENCE_BITS
	sequence := id & int64(maskSequence)

	return []int64{timestamp, nodeId, sequence}
}
