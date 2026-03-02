// Code generated — DO NOT EDIT.

package risk_decision_receipt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb2 "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

var (
	_ = bytes.Equal
	_ = errors.New
	_ = fmt.Sprintf
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
	_ = emptypb.Empty{}
	_ = pb.NewBigIntFromInt
	_ = pb2.AggregationType_AGGREGATION_TYPE_COMMON_PREFIX
	_ = bindings.FilterOptions{}
	_ = evm.FilterLogTriggerRequest{}
	_ = cre.ResponseBufferTooSmall
	_ = rpc.API{}
	_ = json.Unmarshal
	_ = reflect.Bool
)

var RiskDecisionReceiptMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"decisions\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"runId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"decisionHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxPositionUsd\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSlippageBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ttlSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainlinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recorder\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRunCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isDecisionValid\",\"inputs\":[{\"name\":\"runId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordDecision\",\"inputs\":[{\"name\":\"runId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"decisionHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxPositionUsd\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSlippageBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ttlSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainlinkPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recorded\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"runIds\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalApproved\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalDecisions\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalDenied\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DecisionRecorded\",\"inputs\":[{\"name\":\"runId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"decisionHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"recorder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false}]",
}

// Structs

// Contract Method Inputs
type DecisionsInput struct {
	Arg0 [32]byte
}

type IsDecisionValidInput struct {
	RunId [32]byte
}

type RecordDecisionInput struct {
	RunId          [32]byte
	DecisionHash   [32]byte
	Approved       bool
	MaxPositionUsd *big.Int
	MaxSlippageBps *big.Int
	TtlSeconds     *big.Int
	ChainlinkPrice *big.Int
}

type RecordedInput struct {
	Arg0 [32]byte
}

type RunIdsInput struct {
	Arg0 *big.Int
}

// Contract Method Outputs
type DecisionsOutput struct {
	RunId          [32]byte
	DecisionHash   [32]byte
	Approved       bool
	MaxPositionUsd *big.Int
	MaxSlippageBps *big.Int
	TtlSeconds     *big.Int
	ChainlinkPrice *big.Int
	Timestamp      *big.Int
	Recorder       common.Address
}

// Errors

// Events
// The <Event>Topics struct should be used as a filter (for log triggers).
// Note: It is only possible to filter on indexed fields.
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type DecisionRecordedTopics struct {
	RunId        [32]byte
	DecisionHash [32]byte
	Recorder     common.Address
}

type DecisionRecordedDecoded struct {
	RunId        [32]byte
	DecisionHash [32]byte
	Recorder     common.Address
	Approved     bool
}

// Main Binding Type for RiskDecisionReceipt
type RiskDecisionReceipt struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   RiskDecisionReceiptCodec
}

type RiskDecisionReceiptCodec interface {
	EncodeDecisionsMethodCall(in DecisionsInput) ([]byte, error)
	DecodeDecisionsMethodOutput(data []byte) (DecisionsOutput, error)
	EncodeGetRunCountMethodCall() ([]byte, error)
	DecodeGetRunCountMethodOutput(data []byte) (*big.Int, error)
	EncodeIsDecisionValidMethodCall(in IsDecisionValidInput) ([]byte, error)
	DecodeIsDecisionValidMethodOutput(data []byte) (bool, error)
	EncodeRecordDecisionMethodCall(in RecordDecisionInput) ([]byte, error)
	EncodeRecordedMethodCall(in RecordedInput) ([]byte, error)
	DecodeRecordedMethodOutput(data []byte) (bool, error)
	EncodeRunIdsMethodCall(in RunIdsInput) ([]byte, error)
	DecodeRunIdsMethodOutput(data []byte) ([32]byte, error)
	EncodeTotalApprovedMethodCall() ([]byte, error)
	DecodeTotalApprovedMethodOutput(data []byte) (*big.Int, error)
	EncodeTotalDecisionsMethodCall() ([]byte, error)
	DecodeTotalDecisionsMethodOutput(data []byte) (*big.Int, error)
	EncodeTotalDeniedMethodCall() ([]byte, error)
	DecodeTotalDeniedMethodOutput(data []byte) (*big.Int, error)
	DecisionRecordedLogHash() []byte
	EncodeDecisionRecordedTopics(evt abi.Event, values []DecisionRecordedTopics) ([]*evm.TopicValues, error)
	DecodeDecisionRecorded(log *evm.Log) (*DecisionRecordedDecoded, error)
}

func NewRiskDecisionReceipt(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*RiskDecisionReceipt, error) {
	parsed, err := abi.JSON(strings.NewReader(RiskDecisionReceiptMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &RiskDecisionReceipt{
		Address: address,
		Options: options,
		ABI:     &parsed,
		client:  client,
		Codec:   codec,
	}, nil
}

type Codec struct {
	abi *abi.ABI
}

func NewCodec() (RiskDecisionReceiptCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(RiskDecisionReceiptMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeDecisionsMethodCall(in DecisionsInput) ([]byte, error) {
	return c.abi.Pack("decisions", in.Arg0)
}

func (c *Codec) DecodeDecisionsMethodOutput(data []byte) (DecisionsOutput, error) {
	vals, err := c.abi.Methods["decisions"].Outputs.Unpack(data)
	if err != nil {
		return DecisionsOutput{}, err
	}
	if len(vals) != 9 {
		return DecisionsOutput{}, fmt.Errorf("expected 9 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 [32]byte
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 [32]byte
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 bool
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to unmarshal to bool: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 *big.Int
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 *big.Int
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData5, err := json.Marshal(vals[5])
	if err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to marshal ABI result 5: %w", err)
	}

	var result5 *big.Int
	if err := json.Unmarshal(jsonData5, &result5); err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData6, err := json.Marshal(vals[6])
	if err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to marshal ABI result 6: %w", err)
	}

	var result6 *big.Int
	if err := json.Unmarshal(jsonData6, &result6); err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData7, err := json.Marshal(vals[7])
	if err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to marshal ABI result 7: %w", err)
	}

	var result7 *big.Int
	if err := json.Unmarshal(jsonData7, &result7); err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData8, err := json.Marshal(vals[8])
	if err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to marshal ABI result 8: %w", err)
	}

	var result8 common.Address
	if err := json.Unmarshal(jsonData8, &result8); err != nil {
		return DecisionsOutput{}, fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return DecisionsOutput{
		RunId:          result0,
		DecisionHash:   result1,
		Approved:       result2,
		MaxPositionUsd: result3,
		MaxSlippageBps: result4,
		TtlSeconds:     result5,
		ChainlinkPrice: result6,
		Timestamp:      result7,
		Recorder:       result8,
	}, nil
}

func (c *Codec) EncodeGetRunCountMethodCall() ([]byte, error) {
	return c.abi.Pack("getRunCount")
}

func (c *Codec) DecodeGetRunCountMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["getRunCount"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeIsDecisionValidMethodCall(in IsDecisionValidInput) ([]byte, error) {
	return c.abi.Pack("isDecisionValid", in.RunId)
}

func (c *Codec) DecodeIsDecisionValidMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["isDecisionValid"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeRecordDecisionMethodCall(in RecordDecisionInput) ([]byte, error) {
	return c.abi.Pack("recordDecision", in.RunId, in.DecisionHash, in.Approved, in.MaxPositionUsd, in.MaxSlippageBps, in.TtlSeconds, in.ChainlinkPrice)
}

func (c *Codec) EncodeRecordedMethodCall(in RecordedInput) ([]byte, error) {
	return c.abi.Pack("recorded", in.Arg0)
}

func (c *Codec) DecodeRecordedMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["recorded"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeRunIdsMethodCall(in RunIdsInput) ([]byte, error) {
	return c.abi.Pack("runIds", in.Arg0)
}

func (c *Codec) DecodeRunIdsMethodOutput(data []byte) ([32]byte, error) {
	vals, err := c.abi.Methods["runIds"].Outputs.Unpack(data)
	if err != nil {
		return *new([32]byte), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([32]byte), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result [32]byte
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([32]byte), fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTotalApprovedMethodCall() ([]byte, error) {
	return c.abi.Pack("totalApproved")
}

func (c *Codec) DecodeTotalApprovedMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["totalApproved"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTotalDecisionsMethodCall() ([]byte, error) {
	return c.abi.Pack("totalDecisions")
}

func (c *Codec) DecodeTotalDecisionsMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["totalDecisions"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTotalDeniedMethodCall() ([]byte, error) {
	return c.abi.Pack("totalDenied")
}

func (c *Codec) DecodeTotalDeniedMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["totalDenied"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) DecisionRecordedLogHash() []byte {
	return c.abi.Events["DecisionRecorded"].ID.Bytes()
}

func (c *Codec) EncodeDecisionRecordedTopics(
	evt abi.Event,
	values []DecisionRecordedTopics,
) ([]*evm.TopicValues, error) {
	var runIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.RunId).IsZero() {
			runIdRule = append(runIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.RunId)
		if err != nil {
			return nil, err
		}
		runIdRule = append(runIdRule, fieldVal)
	}
	var decisionHashRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.DecisionHash).IsZero() {
			decisionHashRule = append(decisionHashRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.DecisionHash)
		if err != nil {
			return nil, err
		}
		decisionHashRule = append(decisionHashRule, fieldVal)
	}
	var recorderRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Recorder).IsZero() {
			recorderRule = append(recorderRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.Recorder)
		if err != nil {
			return nil, err
		}
		recorderRule = append(recorderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		runIdRule,
		decisionHashRule,
		recorderRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeDecisionRecorded decodes a log into a DecisionRecorded struct.
func (c *Codec) DecodeDecisionRecorded(log *evm.Log) (*DecisionRecordedDecoded, error) {
	event := new(DecisionRecordedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "DecisionRecorded", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["DecisionRecorded"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c RiskDecisionReceipt) Decisions(
	runtime cre.Runtime,
	args DecisionsInput,
	blockNumber *big.Int,
) cre.Promise[DecisionsOutput] {
	calldata, err := c.Codec.EncodeDecisionsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[DecisionsOutput](DecisionsOutput{}, err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (DecisionsOutput, error) {
		return c.Codec.DecodeDecisionsMethodOutput(response.Data)
	})

}

func (c RiskDecisionReceipt) GetRunCount(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeGetRunCountMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeGetRunCountMethodOutput(response.Data)
	})

}

func (c RiskDecisionReceipt) IsDecisionValid(
	runtime cre.Runtime,
	args IsDecisionValidInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeIsDecisionValidMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[bool](*new(bool), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (bool, error) {
		return c.Codec.DecodeIsDecisionValidMethodOutput(response.Data)
	})

}

func (c RiskDecisionReceipt) Recorded(
	runtime cre.Runtime,
	args RecordedInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeRecordedMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[bool](*new(bool), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (bool, error) {
		return c.Codec.DecodeRecordedMethodOutput(response.Data)
	})

}

func (c RiskDecisionReceipt) RunIds(
	runtime cre.Runtime,
	args RunIdsInput,
	blockNumber *big.Int,
) cre.Promise[[32]byte] {
	calldata, err := c.Codec.EncodeRunIdsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[[32]byte](*new([32]byte), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) ([32]byte, error) {
		return c.Codec.DecodeRunIdsMethodOutput(response.Data)
	})

}

func (c RiskDecisionReceipt) TotalApproved(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeTotalApprovedMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeTotalApprovedMethodOutput(response.Data)
	})

}

func (c RiskDecisionReceipt) TotalDecisions(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeTotalDecisionsMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeTotalDecisionsMethodOutput(response.Data)
	})

}

func (c RiskDecisionReceipt) TotalDenied(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeTotalDeniedMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeTotalDeniedMethodOutput(response.Data)
	})

}

func (c RiskDecisionReceipt) WriteReport(
	runtime cre.Runtime,
	report *cre.Report,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  c.Address.Bytes(),
		Report:    report,
		GasConfig: gasConfig,
	})
}

func (c *RiskDecisionReceipt) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	default:
		return nil, errors.New("unknown error selector")
	}
}

// DecisionRecordedTrigger wraps the raw log trigger and provides decoded DecisionRecordedDecoded data
type DecisionRecordedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *RiskDecisionReceipt // Keep reference for decoding
}

// Adapt method that decodes the log into DecisionRecorded data
func (t *DecisionRecordedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[DecisionRecordedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeDecisionRecorded(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DecisionRecorded log: %w", err)
	}

	return &bindings.DecodedLog[DecisionRecordedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *RiskDecisionReceipt) LogTriggerDecisionRecordedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []DecisionRecordedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[DecisionRecordedDecoded]], error) {
	event := c.ABI.Events["DecisionRecorded"]
	topics, err := c.Codec.EncodeDecisionRecordedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for DecisionRecorded: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &DecisionRecordedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *RiskDecisionReceipt) FilterLogsDecisionRecorded(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.DecisionRecordedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}
