// Copyright 2017 Vector Creations Ltd
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

package main

import (
	"github.com/matrix-org/dendrite/internal/setup"
	"github.com/matrix-org/dendrite/userapi"
)

func main() {
	cfg := setup.ParseFlags(false)
	base := setup.NewBaseDendrite(cfg, "UserAPI", true)
	defer base.Close() // nolint: errcheck

	accountDB := base.CreateAccountsDB()
	deviceDB := base.CreateDeviceDB()

	userAPI := userapi.NewInternalAPI(accountDB, deviceDB, cfg.Matrix.ServerName, cfg.Derived.ApplicationServices, base.KeyServerHTTPClient())

	userapi.AddInternalRoutes(base.InternalAPIMux, userAPI)

	base.SetupAndServeHTTP(string(base.Cfg.Bind.UserAPI), string(base.Cfg.Listen.UserAPI))
}