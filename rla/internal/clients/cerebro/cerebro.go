/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cerebro

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/rack"
	identifier "github.com/nvidia/carbide-rest/rla/pkg/common/Identifier"
)

type Client struct {
	client     *resty.Client
	apiKey     string
	graphQLURL string
	cache      *Cache
}

type Config struct {
	URL    string
	APIKey string
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("api key is required")
	}

	if c.URL == "" {
		return fmt.Errorf("url is required")
	}

	return c.verifyAndSetGraphQLURL()
}

func (c *Config) verifyAndSetGraphQLURL() error {
	parsed, err := url.ParseRequestURI(c.URL)
	if err != nil {
		return fmt.Errorf("url is invalid: %v", err)
	}

	trailer := "api/graphql"
	escapedPath := parsed.EscapedPath()
	switch escapedPath {
	case "":
		// Add the full trailer
		c.URL = c.URL + "/" + trailer + "/"
	case "/":
		// Add the trailer without the slash
		c.URL = c.URL + trailer + "/"
	case "/" + trailer:
		// Fulllfill the trailer with the missing slash
		c.URL = c.URL + "/"
	case "/" + trailer + "/":
		// The trailer is already present, nothing to do
	default:
		// The escaped path is not as expected, return error
		return fmt.Errorf("url is not a valid graphQL URL")
	}

	return nil
}

func New(ctx context.Context, conf Config) (*Client, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	client := &Client{
		client: resty.New().
			SetBaseURL(conf.URL).
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json"),
		apiKey:     conf.APIKey,
		graphQLURL: conf.URL,
	}

	if err := client.setupCache(ctx); err != nil {
		return nil, err
	}

	invalidLocations := client.cache.ValidateLocationChain()
	if len(invalidLocations) > 0 {
		return nil, fmt.Errorf("invalid locations: %v", invalidLocations)
	}

	return client, nil
}

const gqlLocationsQuery = `
query ($offset: Int, $limit: Int) {
	locations(offset: $offset, limit: $limit) {
		name
		location_type {
			name
		}
		parent {
			name
		}
	}
}
`
const gqlRackByNameQuery = `
query ($name: [String!]) {
	racks(name: $name, limit: 1) {
		id
		name
		serial
		asset_tag
		location {
			name
			location_type {
				name
			}
			parent {
				name
			}
		}
		devices {
			id
			name
			serial
			position
			role {
				name
			}
			device_bays {
				installed_device {
					id
					name
				}
			}
			device_type {
				model
				manufacturer {
					name
				}
				device_family {
					name
				}
			}
			interfaces (name__re: "bmc") {
				name
				ip_addresses {
					address
				}
				mac_address
			}
		}
	}
}
`

const gqlAllNVLDomainQuery = `
query {
    nvlink_domains (rack__isnull: false) {
        id
        name
    }
}
`

const gqlNVLDomainByIDQuery = `
query ($id: [String!]) {
    nvlink_domains(id: $id) {
        id
        name
        rack {
            id
            name
			asset_tag
			serial
			location {
				name
				location_type {
					name
				}
				parent {
					name
				}
			}
        }
        members {
            id
            name
            serial
			position
			role {
				name
			}
            device_bays {
                installed_device {
                    id
                    name
                }
            }
            device_type{
                model
                manufacturer{
                    name
                }
            }
            interfaces (name__re: "bmc|eth0") {
                name
                ip_addresses {
                    address
                }
                mac_address
            }
        }
    }
}
`

func (c *Client) setupCache(ctx context.Context) error {
	locations, err := c.GetLocations(ctx)
	if err != nil {
		return err
	}

	c.cache = NewCache()
	for _, loc := range locations {
		c.cache.AddLocation(loc.Name, &loc)
	}

	return nil
}

func (c *Client) GetLocations(ctx context.Context) ([]Location, error) {
	total := 0
	offset := 0
	limit := DefaultLocationsLimit
	locations := make([]Location, 0)
	for {
		qreq := Request{
			Query: gqlLocationsQuery,
			Variables: map[string]any{
				"offset": offset,
				"limit":  limit,
			},
		}

		var qresp LocationsResponse

		if err := c.issueRequest(ctx, &qreq, &qresp, false); err != nil {
			return nil, err
		}

		total += len(qresp.Data.Locations)
		if total >= MaxLocationsSupported {
			return nil, fmt.Errorf(
				"max locations supported %d reached",
				MaxLocationsSupported,
			)
		}

		locations = append(locations, qresp.Data.Locations...)

		if len(qresp.Data.Locations) < limit {
			// No more locations to fetch
			break
		}

		offset += limit
	}

	return locations, nil
}

func (c *Client) GetRackByName(
	ctx context.Context,
	name string,
) (*rack.Rack, error) {
	qreq := Request{
		Query: gqlRackByNameQuery,
		Variables: map[string]any{
			"name": name,
		},
	}

	var qresp RackResponse

	if err := c.issueRequest(ctx, &qreq, &qresp, false); err != nil {
		return nil, err
	}

	if len(qresp.Data.Racks) == 0 {
		return nil, fmt.Errorf("rack %s not found", name)
	}

	return newRackConverter(&qresp.Data.Racks[0], c.cache).toRack(), nil
}

func (c *Client) GetAllNVLDomainIdentifiers(
	ctx context.Context,
) ([]identifier.Identifier, error) {
	qreq := Request{
		Query: gqlAllNVLDomainQuery,
	}

	var qresp AllNVLDomainResponse

	if err := c.issueRequest(ctx, &qreq, &qresp, false); err != nil {
		return nil, err
	}

	ids := make([]identifier.Identifier, 0, len(qresp.Data.NVLDomains))
	for _, nvlDomain := range qresp.Data.NVLDomains {
		ids = append(
			ids,
			identifier.Identifier{
				ID:   uuid.MustParse(nvlDomain.ID),
				Name: nvlDomain.Name,
			},
		)
	}

	return ids, nil
}

func (c *Client) GetRackByNVLDomainID(
	ctx context.Context,
	id uuid.UUID,
) ([]*rack.Rack, error) {
	qreq := Request{
		Query: gqlNVLDomainByIDQuery,
		Variables: map[string]any{
			"id": id.String(),
		},
	}

	var qresp NVLDomainResponse

	if err := c.issueRequest(ctx, &qreq, &qresp, false); err != nil {
		return nil, err
	}

	if len(qresp.Data.NVLDomains) == 0 {
		return nil, fmt.Errorf("nvl domain %s not found", id.String())
	}

	racks := make([]*rack.Rack, 0, len(qresp.Data.NVLDomains))
	for _, nvlDomain := range qresp.Data.NVLDomains {
		r := Rack{
			ID:       nvlDomain.Rack.ID,
			Name:     nvlDomain.Rack.Name,
			Serial:   nvlDomain.Rack.Serial,
			Location: nvlDomain.Rack.Location,
			Devices:  nvlDomain.Members,
		}
		racks = append(racks, newRackConverter(&r, c.cache).toRack())
	}

	return racks, nil
}

func (c *Client) issueRequest(
	ctx context.Context,
	qreq any,
	qresp any,
	debugModeOn bool,
) error {
	var qerr CError

	req := c.client.R().SetContext(ctx).SetResult(&qresp).SetError(&qerr).SetBody(qreq)
	req = req.SetHeader("Authorization", fmt.Sprintf("Token %s", c.apiKey))
	if debugModeOn {
		req = req.SetDebug(true).EnableGenerateCurlOnDebug()
	}

	resp, err := req.Post(c.graphQLURL)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("%v (%v)", resp.Status(), qerr.Error())
	}

	return nil
}
