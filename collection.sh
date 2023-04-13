
URL="http://localhost:5001/familytree"
# curl "${URL}/person"

# create people Sonny, Mike, Martin, Phoebe, Anastasia, Ellen, Ursula, Oprah, Eric, Ariel, Dunny, Bruce, Jacqueline, Melody 
curl -X POST \
  -H "Content-Type: application/json" \
  -d @- \
"${URL}/people" <<EOF
[	{		"name": "Sonny"	},	{		"name": "Mike"	},	{		"name": "Martin"	},	{		"name": "Phoebe"	},	{		"name": "Anastasia"	},	{		"name": "Ellen"	},	{		"name": "Ursula"	},	{		"name": "Oprah"	},	{		"name": "Eric"	},	{		"name": "Ariel"	},	{		"name": "Dunny"	},	{		"name": "Bruce"	},	{		"name": "Jacqueline"	},	{		"name": "Melody"	}]
EOF

curl -X POST \
  -H "Content-Type: application/json" \
  -d @- \
  "http://localhost:5001/familytree/relationships" <<EOF
  [ {
    "children": "e4c04e63-4982-4357-8c26-dc277c76779a",
    "parent": "393114f0-704b-4214-936a-97a804a14e71"
  }]
EOF
