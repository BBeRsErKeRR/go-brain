package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate/inventory"
)

func main() {
	input := []byte(`{
	"_meta": {
		"hostvars": {
			"facthost0": {
				"ansible_connection": "local",
				"ansible_host": "1270.0.0.1"
			},
			"facthost1": {
				"ansible_connection": "local",
				"ansible_host": "1270.0.0.1"
			},
			"facthost2": {
				"ansible_connection": "local",
				"ansible_host": "1270.0.0.1"
			},
			"facthost3": {
				"ansible_connection": "local",
				"ansible_host": "1270.0.0.1"
			},
			"localhost": {
				"ansible_connection": "local",
				"ansible_ssh_host": "127.0.0.1"
			},
			"testhost": {
				"a": 1,
				"ansible_connection": "local",
				"ansible_ssh_host": "127.0.0.1",
				"b": 2,
				"c": 3,
				"d": 4,
				"defaults_file_var_role3": "overridden from inventory",
				"role_var_beats_inventory": "should_not_see_this",
				"test_hash": {
					"host_vars_testhost": "this is in host_vars/testhost"
				}
			},
			"testhost2": {
				"ansible_connection": "local",
				"ansible_ssh_host": "127.0.0.1"
			},
			"testhost3": {
				"ansible_ssh_host": "127.0.0.3"
			},
			"testhost4": {
				"ansible_ssh_host": "127.0.0.4"
			}
		}
	},
	"all": {
		"children": [
			"ungrouped",
			"inven_overridehosts",
			"arbitrary_grandparent",
			"amazon"
		],
		"hosts": [],
		"vars": {
			"a": 999,
			"b": 998,
			"c": 997,
			"d": 996,
			"dos": 2,
			"etest": "from group_vars",
			"extra_var_override": "FROM_INVENTORY",
			"inven_var": "inventory_var",
			"inventory_beats_default": "narf",
			"test_bare": true,
			"test_bare_nested_bad": "{{test_bare_var}} == 321",
			"test_bare_nested_good": "{{test_bare_var}} == 123",
			"test_bare_var": 123,
			"test_hash": {
				"group_vars_all": "this is in group_vars/all"
			},
			"tres": 3,
			"unicode_host_var": "Caf\u00e9E\u00f1yei",
			"uno": 1
		}
	},
	"amazon": {
		"children": [],
		"hosts": [
			"localhost"
		],
		"vars": {
			"ec2_region": "us-east-1",
			"ec2_url": "ec2.amazonaws.com"
		}
	},
	"arbitrary_grandparent": {
		"children": [
			"arbitrary_parent"
		],
		"hosts": [],
		"vars": {
			"grandparent_var": 2000,
			"groups_tree_var": 3000,
			"overridden_in_parent": 2000
		}
	},
	"arbitrary_parent": {
		"children": [
			"local"
		],
		"hosts": [],
		"vars": {
			"groups_tree_var": 4000,
			"overridden_in_parent": 1000
		}
	},
	"inven_overridehosts": {
		"children": [],
		"hosts": [
			"invenoverride"
		],
		"vars": {
			"foo": "foo",
			"var_dir": "vars"
		}
	},
	"local": {
		"children": [],
		"hosts": [
			"testhost",
			"testhost2",
			"testhost3"
		],
		"vars": {
			"groups_tree_var": 5000,
			"hash_test": {
				"group_vars_local": "this is in group_vars/local"
			},
			"parent_var": 6000,
			"tres": "three"
		}
	},
	"ungrouped": {
		"children": [],
		"hosts": [],
		"vars": {}
	}
}`)
	inv, err := inventory.NewInventoryFromJson(input)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(inv)
	res, err := json.Marshal(inv)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(res))
}
