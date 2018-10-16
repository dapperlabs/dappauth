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

// ERC725CoreABI is the input ABI used to generate the binding from.
const ERC725CoreABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_key\",\"type\":\"bytes32\"},{\"name\":\"_purpose\",\"type\":\"uint256\"}],\"name\":\"keyHasPurpose\",\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ERC725Core is an auto generated Go binding around an Ethereum contract.
type ERC725Core struct {
	ERC725CoreCaller     // Read-only binding to the contract
	ERC725CoreTransactor // Write-only binding to the contract
	ERC725CoreFilterer   // Log filterer for contract events
}

// ERC725CoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC725CoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC725CoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC725CoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC725CoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC725CoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC725CoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC725CoreSession struct {
	Contract     *ERC725Core       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC725CoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC725CoreCallerSession struct {
	Contract *ERC725CoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ERC725CoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC725CoreTransactorSession struct {
	Contract     *ERC725CoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ERC725CoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC725CoreRaw struct {
	Contract *ERC725Core // Generic contract binding to access the raw methods on
}

// ERC725CoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC725CoreCallerRaw struct {
	Contract *ERC725CoreCaller // Generic read-only contract binding to access the raw methods on
}

// ERC725CoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC725CoreTransactorRaw struct {
	Contract *ERC725CoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC725Core creates a new instance of ERC725Core, bound to a specific deployed contract.
func NewERC725Core(address common.Address, backend bind.ContractBackend) (*ERC725Core, error) {
	contract, err := bindERC725Core(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC725Core{ERC725CoreCaller: ERC725CoreCaller{contract: contract}, ERC725CoreTransactor: ERC725CoreTransactor{contract: contract}, ERC725CoreFilterer: ERC725CoreFilterer{contract: contract}}, nil
}

// NewERC725CoreCaller creates a new read-only instance of ERC725Core, bound to a specific deployed contract.
func NewERC725CoreCaller(address common.Address, caller bind.ContractCaller) (*ERC725CoreCaller, error) {
	contract, err := bindERC725Core(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC725CoreCaller{contract: contract}, nil
}

// NewERC725CoreTransactor creates a new write-only instance of ERC725Core, bound to a specific deployed contract.
func NewERC725CoreTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC725CoreTransactor, error) {
	contract, err := bindERC725Core(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC725CoreTransactor{contract: contract}, nil
}

// NewERC725CoreFilterer creates a new log filterer instance of ERC725Core, bound to a specific deployed contract.
func NewERC725CoreFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC725CoreFilterer, error) {
	contract, err := bindERC725Core(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC725CoreFilterer{contract: contract}, nil
}

// bindERC725Core binds a generic wrapper to an already deployed contract.
func bindERC725Core(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC725CoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC725Core *ERC725CoreRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC725Core.Contract.ERC725CoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC725Core *ERC725CoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC725Core.Contract.ERC725CoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC725Core *ERC725CoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC725Core.Contract.ERC725CoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC725Core *ERC725CoreCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC725Core.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC725Core *ERC725CoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC725Core.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC725Core *ERC725CoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC725Core.Contract.contract.Transact(opts, method, params...)
}

// KeyHasPurpose is a free data retrieval call binding the contract method 0xd202158d.
//
// Solidity: function keyHasPurpose(_key bytes32, _purpose uint256) constant returns(exists bool)
func (_ERC725Core *ERC725CoreCaller) KeyHasPurpose(opts *bind.CallOpts, _key [32]byte, _purpose *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC725Core.contract.Call(opts, out, "keyHasPurpose", _key, _purpose)
	return *ret0, err
}

// KeyHasPurpose is a free data retrieval call binding the contract method 0xd202158d.
//
// Solidity: function keyHasPurpose(_key bytes32, _purpose uint256) constant returns(exists bool)
func (_ERC725Core *ERC725CoreSession) KeyHasPurpose(_key [32]byte, _purpose *big.Int) (bool, error) {
	return _ERC725Core.Contract.KeyHasPurpose(&_ERC725Core.CallOpts, _key, _purpose)
}

// KeyHasPurpose is a free data retrieval call binding the contract method 0xd202158d.
//
// Solidity: function keyHasPurpose(_key bytes32, _purpose uint256) constant returns(exists bool)
func (_ERC725Core *ERC725CoreCallerSession) KeyHasPurpose(_key [32]byte, _purpose *big.Int) (bool, error) {
	return _ERC725Core.Contract.KeyHasPurpose(&_ERC725Core.CallOpts, _key, _purpose)
}
