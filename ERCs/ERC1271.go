// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ERCs

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ERC1271ABI is the input ABI used to generate the binding from.
const ERC1271ABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_data\",\"type\":\"bytes\"},{\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"name\":\"magicValue\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ERC1271 is an auto generated Go binding around an Ethereum contract.
type ERC1271 struct {
	ERC1271Caller     // Read-only binding to the contract
	ERC1271Transactor // Write-only binding to the contract
	ERC1271Filterer   // Log filterer for contract events
}

// ERC1271Caller is an auto generated read-only Go binding around an Ethereum contract.
type ERC1271Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1271Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC1271Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1271Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC1271Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1271Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC1271Session struct {
	Contract     *ERC1271          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC1271CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC1271CallerSession struct {
	Contract *ERC1271Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ERC1271TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC1271TransactorSession struct {
	Contract     *ERC1271Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ERC1271Raw is an auto generated low-level Go binding around an Ethereum contract.
type ERC1271Raw struct {
	Contract *ERC1271 // Generic contract binding to access the raw methods on
}

// ERC1271CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC1271CallerRaw struct {
	Contract *ERC1271Caller // Generic read-only contract binding to access the raw methods on
}

// ERC1271TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC1271TransactorRaw struct {
	Contract *ERC1271Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1271 creates a new instance of ERC1271, bound to a specific deployed contract.
func NewERC1271(address common.Address, backend bind.ContractBackend) (*ERC1271, error) {
	contract, err := bindERC1271(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1271{ERC1271Caller: ERC1271Caller{contract: contract}, ERC1271Transactor: ERC1271Transactor{contract: contract}, ERC1271Filterer: ERC1271Filterer{contract: contract}}, nil
}

// NewERC1271Caller creates a new read-only instance of ERC1271, bound to a specific deployed contract.
func NewERC1271Caller(address common.Address, caller bind.ContractCaller) (*ERC1271Caller, error) {
	contract, err := bindERC1271(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1271Caller{contract: contract}, nil
}

// NewERC1271Transactor creates a new write-only instance of ERC1271, bound to a specific deployed contract.
func NewERC1271Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC1271Transactor, error) {
	contract, err := bindERC1271(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1271Transactor{contract: contract}, nil
}

// NewERC1271Filterer creates a new log filterer instance of ERC1271, bound to a specific deployed contract.
func NewERC1271Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC1271Filterer, error) {
	contract, err := bindERC1271(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1271Filterer{contract: contract}, nil
}

// bindERC1271 binds a generic wrapper to an already deployed contract.
func bindERC1271(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC1271ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1271 *ERC1271Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC1271.Contract.ERC1271Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1271 *ERC1271Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1271.Contract.ERC1271Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1271 *ERC1271Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1271.Contract.ERC1271Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1271 *ERC1271CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC1271.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1271 *ERC1271TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1271.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1271 *ERC1271TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1271.Contract.contract.Transact(opts, method, params...)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x20c13b0b.
//
// Solidity: function isValidSignature(_data bytes, _signature bytes) constant returns(magicValue bytes4)
func (_ERC1271 *ERC1271Caller) IsValidSignature(opts *bind.CallOpts, _data []byte, _signature []byte) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _ERC1271.contract.Call(opts, out, "isValidSignature", _data, _signature)
	return *ret0, err
}

// IsValidSignature is a free data retrieval call binding the contract method 0x20c13b0b.
//
// Solidity: function isValidSignature(_data bytes, _signature bytes) constant returns(magicValue bytes4)
func (_ERC1271 *ERC1271Session) IsValidSignature(_data []byte, _signature []byte) ([4]byte, error) {
	return _ERC1271.Contract.IsValidSignature(&_ERC1271.CallOpts, _data, _signature)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x20c13b0b.
//
// Solidity: function isValidSignature(_data bytes, _signature bytes) constant returns(magicValue bytes4)
func (_ERC1271 *ERC1271CallerSession) IsValidSignature(_data []byte, _signature []byte) ([4]byte, error) {
	return _ERC1271.Contract.IsValidSignature(&_ERC1271.CallOpts, _data, _signature)
}
