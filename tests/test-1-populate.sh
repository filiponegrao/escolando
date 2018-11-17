# Authenticacao

loginResponse="$(curl -X POST http://localhost:8080/login -H 'Content-Type:application/json' -d '{"username": "filiponegrao@gmail.com", "password": "teste"}')"
auth=(${loginResponse#*token})
auth=(${auth:3})
len=(${#auth})
auth="Authorization: Bearer ${auth:0:len-2}"

# Cria a insittuicao
resp="$(curl -X POST localhost:8080/institutions -H "${auth}" -H 'Content-Type: application/json' -d '{"name":"Escolando - educação do futuro", "email": "contato@escolinha.com", "addressStreet": "Rua São Clemente, Botafogo RJ", "addressNumber": 10000, "addressPostal": "26722-901", "owner": {"id": 1}}')"

error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error} " ao criar instituicao"
	exit
fi

# Recupera o id da resposta
institutionId=$(echo ${resp} | cut -d ',' -f 1)
institutionId=$(echo ${institutionId} | cut -d ':' -f 2) 


# # Cria o segmento Ensino Fundamental
segment='{"name":"Ensino Fundamental", "institution": {"id":'${institutionId}'}}'
request="curl -X POST localhost:8080/segments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${segment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then 
	clear
	echo ${error}
	echo ${request}
	exit
fi
segment1Id=$(echo ${resp} | cut -d ',' -f 1)
segment1Id=$(echo ${segment1Id} | cut -d ':' -f 2) 


# # Cria o segmento Ensino Médio
segment='{"name":"Ensino Médio", "institution": {"id":'${institutionId}'}}'
request="curl -X POST localhost:8080/segments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${segment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then 
	clear
	echo ${error}
	echo ${request}
	exit
fi
segment2Id=$(echo ${resp} | cut -d ',' -f 1)
segment2Id=$(echo ${segment2Id} | cut -d ':' -f 2) 

# Cria a primeira serie
grade='{"name": "1ª Série", "segment": {"id": '${segment1Id}'}}'
request="curl -X POST localhost:8080/schoolGrades -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${grade})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then 
	clear
	echo ${error}
	echo ${request}
	exit
fi
grade1Id=$(echo ${resp} | cut -d ',' -f 1)
grade1Id=$(echo ${grade1Id} | cut -d ':' -f 2) 

# Cria a segunda serie
grade='{"name": "2ª Série", "segment": {"id": '${segment1Id}'}}'
request="curl -X POST localhost:8080/schoolGrades -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${grade})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
grade2Id=$(echo ${resp} | cut -d ',' -f 1)
grade2Id=$(echo ${grade2Id} | cut -d ':' -f 2) 

# Cria a teceira serie
grade='{"name": "3ª Série", "segment": {"id": '${segment1Id}'}}'
request="curl -X POST localhost:8080/schoolGrades -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${grade})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
grade3Id=$(echo ${resp} | cut -d ',' -f 1)
grade3Id=$(echo ${grade3Id} | cut -d ':' -f 2) 

# Cria o primeiro ano
grade='{"name": "1º Ano", "segment": {"id": '${segment2Id}'}}'
request="curl -X POST localhost:8080/schoolGrades -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${grade})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
grade4Id=$(echo ${resp} | cut -d ',' -f 1)
grade4Id=$(echo ${grade4Id} | cut -d ':' -f 2) 


# Cria o segundo ano
grade='{"name": "2º Ano", "segment": {"id": '${segment2Id}'}}'
request="curl -X POST localhost:8080/schoolGrades -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${grade})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
grade5Id=$(echo ${resp} | cut -d ',' -f 1)
grade5Id=$(echo ${grade5Id} | cut -d ':' -f 2) 

# Cria o terceiro ano
grade='{"name": "3º Ano", "segment": {"id": '${segment2Id}'}}'
request="curl -X POST localhost:8080/schoolGrades -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${grade})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
grade6Id=$(echo ${resp} | cut -d ',' -f 1)
grade6Id=$(echo ${grade6Id} | cut -d ':' -f 2) 

### Encarregados ###

# Cria o diretor

incharge='{"name": "Alcino", "email": "alcino@escolando.com", "role": {"id": 1}, "institution":{"id": '${institutionId}'} }'
request="curl -X POST localhost:8080/inChargeUser -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${incharge})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
incharge1Id=$(echo ${resp} | cut -d ',' -f 1)
incharge1Id=$(echo ${incharge1Id} | cut -d ':' -f 2) 

# Cria a coordenadora

incharge='{"name": "Barbara Bara", "email": "barbara@escolando.com", "role": {"id": 2}, "institution":{"id": '${institutionId}'} }'
request="curl -X POST localhost:8080/inChargeUser -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${incharge})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
incharge2Id=$(echo ${resp} | cut -d ',' -f 1)
incharge2Id=$(echo ${incharge2Id} | cut -d ':' -f 2) 

# Cria a secretaria
incharge='{"name": "Clarinéia", "email": "clarineia@escolando.com", "role": {"id": 3}, "institution":{"id": '${institutionId}'} }'
request="curl -X POST localhost:8080/inChargeUser -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${incharge})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
incharge3Id=$(echo ${resp} | cut -d ',' -f 1)
incharge3Id=$(echo ${incharge3Id} | cut -d ':' -f 2) 

# Cria o professor
incharge='{"name": "Damastor", "email": "damastor@escolando.com", "role": {"id": 4}, "institution":{"id": '${institutionId}'} }'
request="curl -X POST localhost:8080/inChargeUser -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${incharge})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
incharge4Id=$(echo ${resp} | cut -d ',' -f 1)
incharge4Id=$(echo ${incharge4Id} | cut -d ':' -f 2) 

# Cria a professora
incharge='{"name": "Edineia", "email": "edineia@escolando.com", "role": {"id": 4}, "institution":{"id": '${institutionId}'} }'
request="curl -X POST localhost:8080/inChargeUser -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${incharge})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
incharge5Id=$(echo ${resp} | cut -d ',' -f 1)
incharge5Id=$(echo ${incharge5Id} | cut -d ':' -f 2) 

### Turmas ###

# Turma A
class='{"name": "Turma A", "capacity": 10, "schoolGrade": {"id": '${grade1Id}'}, "inCharge": {"id": '${incharge2Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class1Id=$(echo ${resp} | cut -d ',' -f 1)
class1Id=$(echo ${class1Id} | cut -d ':' -f 2) 

# Turma B
class='{"name": "Turma B", "capacity": 10, "schoolGrade": {"id": '${grade1Id}'}, "inCharge": {"id": '${incharge2Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class2Id=$(echo ${resp} | cut -d ',' -f 1)
class2Id=$(echo ${class2Id} | cut -d ':' -f 2) 

# Turma C
class='{"name": "Turma C", "capacity": 10, "schoolGrade": {"id": '${grade1Id}'}, "inCharge": {"id": '${incharge2Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class3Id=$(echo ${resp} | cut -d ',' -f 1)
class3Id=$(echo ${class3Id} | cut -d ':' -f 2) 

# Turma 1
class='{"name": "Turma 1", "capacity": 10, "schoolGrade": {"id": '${grade2Id}'}, "inCharge": {"id": '${incharge2Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class4Id=$(echo ${resp} | cut -d ',' -f 1)
class4Id=$(echo ${class4Id} | cut -d ':' -f 2) 

# Turma 2
class='{"name": "Turma 2", "capacity": 10, "schoolGrade": {"id": '${grade2Id}'}, "inCharge": {"id": '${incharge2Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class5Id=$(echo ${resp} | cut -d ',' -f 1)
class5Id=$(echo ${class5Id} | cut -d ':' -f 2)

# Turma X
class='{"name": "Turma X", "capacity": 10, "schoolGrade": {"id": '${grade3Id}'}, "inCharge": {"id": '${incharge2Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class6Id=$(echo ${resp} | cut -d ',' -f 1)
class6Id=$(echo ${class6Id} | cut -d ':' -f 2)

# Turma A1
class='{"name": "Turma A1", "capacity": 10, "schoolGrade": {"id": '${grade4Id}'}, "inCharge": {"id": '${incharge4Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class7Id=$(echo ${resp} | cut -d ',' -f 1)
class7Id=$(echo ${class7Id} | cut -d ':' -f 2)

# Turma A2
class='{"name": "Turma A2", "capacity": 10, "schoolGrade": {"id": '${grade4Id}'}, "inCharge": {"id": '${incharge4Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class8Id=$(echo ${resp} | cut -d ',' -f 1)
class8Id=$(echo ${class8Id} | cut -d ':' -f 2)

# Turma B1
class='{"name": "Turma B1", "capacity": 10, "schoolGrade": {"id": '${grade5Id}'}, "inCharge": {"id": '${incharge4Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class9Id=$(echo ${resp} | cut -d ',' -f 1)
class9Id=$(echo ${class9Id} | cut -d ':' -f 2)

# Turma C1
class='{"name": "Turma C1", "capacity": 10, "schoolGrade": {"id": '${grade6Id}'}, "inCharge": {"id": '${incharge4Id}'}}'
request="curl -X POST localhost:8080/classes -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${class})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
class10Id=$(echo ${resp} | cut -d ',' -f 1)
class10Id=$(echo ${class10Id} | cut -d ':' -f 2)




echo "FIM DE ROTINA"