// Copyright 2024 Cysic Labs
// This file is part of Cysic Labs' Cysicmint library.
//
// The Cysicmint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Cysicmint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Cysicmint library. If not, see https://github.com/cysic-labs/cysic-network/blob/main/LICENSE
package types

import (
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// EVMTxIndexer defines the interface of custom eth tx indexer.
type EVMTxIndexer interface {
	// LastIndexedBlock returns -1 if indexer db is empty
	LastIndexedBlock() (int64, error)
	IndexBlock(*tmtypes.Block, []*abci.ResponseDeliverTx) error

	// GetByTxHash returns nil if tx not found.
	GetByTxHash(common.Hash) (*TxResult, error)
	// GetByBlockAndIndex returns nil if tx not found.
	GetByBlockAndIndex(int64, int32) (*TxResult, error)
}
