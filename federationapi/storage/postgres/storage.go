// Copyright 2017-2018 New Vector Ltd
// Copyright 2019-2020 The Matrix.org Foundation C.I.C.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import (
	"database/sql"
	"fmt"

	"github.com/matrix-org/dendrite/federationapi/storage/postgres/deltas"
	"github.com/matrix-org/dendrite/federationapi/storage/shared"
	"github.com/matrix-org/dendrite/internal/caching"
	"github.com/matrix-org/dendrite/internal/sqlutil"
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/gomatrixserverlib"
)

// Database stores information needed by the federation sender
type Database struct {
	shared.Database
	db     *sql.DB
	writer sqlutil.Writer
}

// NewDatabase opens a new database
func NewDatabase(dbProperties *config.DatabaseOptions, cache caching.FederationCache, serverName gomatrixserverlib.ServerName) (*Database, error) {
	var d Database
	var err error
	if d.db, err = sqlutil.Open(dbProperties); err != nil {
		return nil, err
	}
	d.writer = sqlutil.NewDummyWriter()
	joinedHosts, err := NewPostgresJoinedHostsTable(d.db)
	if err != nil {
		return nil, err
	}
	queuePDUs, err := NewPostgresQueuePDUsTable(d.db)
	if err != nil {
		return nil, err
	}
	queueEDUs, err := NewPostgresQueueEDUsTable(d.db)
	if err != nil {
		return nil, err
	}
	queueJSON, err := NewPostgresQueueJSONTable(d.db)
	if err != nil {
		return nil, err
	}
	blacklist, err := NewPostgresBlacklistTable(d.db)
	if err != nil {
		return nil, err
	}
	inboundPeeks, err := NewPostgresInboundPeeksTable(d.db)
	if err != nil {
		return nil, err
	}
	outboundPeeks, err := NewPostgresOutboundPeeksTable(d.db)
	if err != nil {
		return nil, err
	}
	notaryJSON, err := NewPostgresNotaryServerKeysTable(d.db)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresNotaryServerKeysTable: %s", err)
	}
	notaryMetadata, err := NewPostgresNotaryServerKeysMetadataTable(d.db)
	if err != nil {
		return nil, fmt.Errorf("NewPostgresNotaryServerKeysMetadataTable: %s", err)
	}
	serverSigningKeys, err := NewPostgresServerSigningKeysTable(d.db)
	if err != nil {
		return nil, err
	}
	m := sqlutil.NewMigrations()
	deltas.LoadRemoveRoomsTable(m)
	if err = m.RunDeltas(d.db, dbProperties); err != nil {
		return nil, err
	}
	d.Database = shared.Database{
		DB:                       d.db,
		ServerName:               serverName,
		Cache:                    cache,
		Writer:                   d.writer,
		FederationJoinedHosts:    joinedHosts,
		FederationQueuePDUs:      queuePDUs,
		FederationQueueEDUs:      queueEDUs,
		FederationQueueJSON:      queueJSON,
		FederationBlacklist:      blacklist,
		FederationInboundPeeks:   inboundPeeks,
		FederationOutboundPeeks:  outboundPeeks,
		NotaryServerKeysJSON:     notaryJSON,
		NotaryServerKeysMetadata: notaryMetadata,
		ServerSigningKeys:        serverSigningKeys,
	}
	return &d, nil
}
