/*
Copyright 2015 Container Solutions

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cobblerclient

// TransactionBegin starts a per-token transaction. Subsequent modifications
// made with the same token are isolated until TransactionCommit or
// TransactionAbort. Added in Cobbler 4.0.0.
func (c *Client) TransactionBegin() error {
	_, err := c.Call("transaction_begin", c.Token)
	return err
}

// TransactionCommit atomically commits all modifications queued in the current
// transaction. Added in Cobbler 4.0.0.
func (c *Client) TransactionCommit() error {
	_, err := c.Call("transaction_commit", c.Token)
	return err
}

// TransactionAbort rolls back all modifications queued in the current
// transaction. Added in Cobbler 4.0.0.
func (c *Client) TransactionAbort() error {
	_, err := c.Call("transaction_abort", c.Token)
	return err
}
