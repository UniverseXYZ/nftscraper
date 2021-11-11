package contract

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"github.com/universexyz/nftscraper/contract/erc165"
)

var (
	erc165CallIfaceID        = common.Hex2Bytes("0x01ffc9a701ffc9a700000000000000000000000000000000000000000000000000000000")
	erc165CallInvalidIfaceID = common.Hex2Bytes("0x01ffc9a7ffffffff00000000000000000000000000000000000000000000000000000000")
)

// ERC165Supports check if the given contract implements ERC165
func ERC165Supports(ctx context.Context, caller bind.ContractCaller, contractAddr common.Address) (bool, error) {
	retData, err := caller.CallContract(ctx, ethereum.CallMsg{
		Data: erc165CallIfaceID,
	}, nil)

	if err != nil {
		return false, errors.Wrap(err, "1-"+err.Error())
	}

	if len(retData) == 0 {
		return false, nil
	}

	var res bool
	if err := rlp.DecodeBytes(retData, &res); err != nil {
		return false, errors.Wrap(err, "2-"+common.Bytes2Hex(retData)+"-"+err.Error())
	}

	if !res {
		return false, nil
	}

	retData, err = caller.CallContract(ctx, ethereum.CallMsg{
		Data: erc165CallInvalidIfaceID,
	}, nil)

	if err != nil {
		return false, errors.Wrap(err, "3-"+err.Error())
	}

	if len(retData) == 0 {
		return false, nil
	}

	if err := rlp.DecodeBytes(retData, &res); err != nil {
		return false, errors.Wrap(err, "4-"+err.Error())
	}

	if !res {
		return false, nil
	}

	return true, nil
}

func ERC165Implements(ctx context.Context, caller bind.ContractCaller, contractAddr common.Address, ifaceID [4]byte) (bool, error) {
	retData, err := caller.CallContract(ctx, ethereum.CallMsg{
		Data: erc165CallIfaceID,
	}, nil)

	if err != nil {
		return false, errors.Wrap(err, "1-"+err.Error())
	}

	var res bool
	if err := rlp.DecodeBytes(retData, &res); err != nil {
		return false, errors.Wrap(err, "2-"+common.Bytes2Hex(retData)+"-"+err.Error())
	}

	if !res {
		return false, errors.New("5")
	}

	retData, err = caller.CallContract(ctx, ethereum.CallMsg{
		Data: erc165CallInvalidIfaceID,
	}, nil)

	if err != nil {
		return false, errors.Wrap(err, "3-"+err.Error())
	}

	if err := rlp.DecodeBytes(retData, &res); err != nil {
		return false, errors.Wrap(err, "4-"+err.Error())
	}

	if !res {
		return false, errors.New("6")
	}

	contract, err := erc165.NewERC165Caller(contractAddr, caller)
	if err != nil {
		return false, errors.Wrap(err, "7-"+err.Error())
	}

	ret, err := contract.SupportsInterface(&bind.CallOpts{Context: ctx}, ifaceID)
	if err != nil {
		return false, errors.Wrap(err, "8-"+err.Error())
	}

	return ret, nil
}
