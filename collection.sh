
URL="http://localhost:5001/familytree"
# curl "${URL}/person"

# create people Sonny, Mike, Martin, Phoebe, Anastasia, Ellen, Ursula, Oprah, Eric, Ariel, Dunny, Bruce, Jacqueline, Melody 
curl -X POST \
  -H "Content-Type: application/json" \
  -d @- \
"${URL}/people" <<EOF
[	{		"name": "Sonny"	},	{		"name": "Mike"	},	{		"name": "Martin"	},	{		"name": "Phoebe"	},	{		"name": "Anastasia"	},	{		"name": "Ellen"	},	{		"name": "Ursula"	},	{		"name": "Oprah"	},	{		"name": "Eric"	},	{		"name": "Ariel"	},	{		"name": "Dunny"	},	{		"name": "Bruce"	},	{		"name": "Jacqueline"	},	{		"name": "Melody"	}]
EOF



