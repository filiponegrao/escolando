# Authenticacao

loginResponse="$(curl -X POST http://localhost:8080/login -H 'Content-Type:application/json' -d '{"username": "filiponegrao@gmail.com", "password": "teste"}')"
auth=(${loginResponse#*token})
auth=(${auth:3})
len=(${#auth})
auth="Authorization: Bearer ${auth:0:len-2}"

# Cria a insittuicao
resp="$(curl -X POST localhost:8080/institutions -H "${auth}" -H 'Content-Type: application/json' -d '{"name":"Escolando - Small Data", "email": "contato@escolinha.com", "addressStreet": "Rua São Clemente, Botafogo RJ", "addressNumber": 10000, "addressPostal": "26722-901", "owner": {"id": 1}}')"

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

### PARENTES E ALUNOS ###

## Cria Ana e aninha
relation='{"student": {"name":"Aninha", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Ana", "email": "ana@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://images.unsplash.com/photo-1529934901952-dbf92b318361?ixlib=rb-0.3.5&ixid=eyJhcHBfaWQiOjEyMDd9&s=858cd99a43297eaa15764c927bf5faa8&auto=format&fit=crop&w=934&q=80", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 1}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user1Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student1Id=${studentId}

## Cria Bia e biazinha
relation='{"student": {"name":"Biazinha", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Bia", "email": "bia@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 1}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user2Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student2Id=${studentId}

## Cria Carlos e carlito
relation='{"student": {"name":"Carlito", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Carlos", "email": "carlos@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user3Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student3Id=${studentId}

## Cria Diego e dieguito
relation='{"student": {"name":"Dieguito", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Diego", "email": "diego@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user4Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student4Id=${studentId}

## Cria Ester e esterzinha
relation='{"student": {"name":"Esterzinha", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Ester", "email": "ester@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 1}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user5Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student5Id=${studentId}

## Cria Fernanda e fernandinha
relation='{"student": {"name":"Fernandinha", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Fernanda", "email": "fernanda@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://scontent.fsdu5-1.fna.fbcdn.net/v/t1.0-9/21462447_1388065451249309_1504736756059084017_n.jpg?_nc_cat=110&_nc_ht=scontent.fsdu5-1.fna&oh=32561583b027bf7e41608f7696c2e7c4&oe=5C6FAC00", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 1}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user6Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student6Id=${studentId}

## Cria Gabriel e gabrielzinho
relation='{"student": {"name":"Gabrielzinho", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Gabriel", "email": "gabriel@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user7Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student7Id=${studentId}

## Cria Hugo e huguinho
relation='{"student": {"name":"Huguinho", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Hugo", "email": "hugo@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user8Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student8Id=${studentId}

## Cria Ian e Ianzinho
relation='{"student": {"name":"Ianzinho", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Ian", "email": "ian@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user9Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student9Id=${studentId}

## Cria Joao e Joaozinho
relation='{"student": {"name":"Joaozinho", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Joao", "email": "joao@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user10Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student10Id=${studentId}

## Cria Karina e karininha
relation='{"student": {"name":"Karininha", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Karina", "email": "karina@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user11Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student11Id=${studentId}

## Cria Larissa e larissinha
relation='{"student": {"name":"Larissinha", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Larissa", "email": "larissa@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user12Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student12Id=${studentId}

## Cria Maria e Mariazinha
relation='{"student": {"name":"Mariazinha", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Maria", "email": "maria@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user13Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student13Id=${studentId}

## Cria Nathan e Nathanzinho
relation='{"student": {"name":"Nathanzinho", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Nathan", "email": "nathan@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user14Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student14Id=${studentId}

## Cria Oswaldo e owswalinho
relation='{"student": {"name":"Oswaldinho", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Oswaldo", "email": "oswaldo@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user15Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student15Id=${studentId}

## Cria Pedro e pedrinho
relation='{"student": {"name":"Pedrinho", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Pedro", "email": "pedro@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user16Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student16Id=${studentId}

## Cria Quézia e quézinha
relation='{"student": {"name":"Quézinha", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Quézia", "email": "quézia@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user17Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student17Id=${studentId}

## Cria Rodrigo e rodriguinho
relation='{"student": {"name":"Rodriguinho", "institution": {"id": '${institutionId}'}}, "parent": {"name": "Rodrigo", "email": "oswaldo@escolando.com", "phone1": "(21)99956-8957", "profileImageUrl": "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973461_1280.png", "institution": {"id": '${institutionId}'}}, "kinship": {"id": 2}}'
request="curl -X POST localhost:8080/userParentAndStudent -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${relation})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
userId=$(echo ${resp} | awk -F 'parent":' '{print $2}')
userId=$(echo ${userId} | awk -F 'userId":' '{print $2}')
userId=$(echo ${userId} | cut -d ',' -f 1)
user18Id=${userId}
studentId=$(echo ${resp} | awk -F 'student":' '{print $2}')
studentId=$(echo ${studentId} | awk -F 'id":' '{print $2}')
studentId=$(echo ${studentId} | cut -d ',' -f 1)
student18Id=${studentId}

### Matricula de alunos nas turmas ###

enrollment='{"student": {"id": '${student1Id}'}, "class": {"id": '${class1Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)

enrollment='{"student": {"id": '${student2Id}'}, "class": {"id": '${class2Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student3Id}'}, "class": {"id": '${class3Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)

enrollment='{"student": {"id": '${student4Id}'}, "class": {"id": '${class4Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)

enrollment='{"student": {"id": '${student5Id}'}, "class": {"id": '${class5Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student6Id}'}, "class": {"id": '${class6Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student7Id}'}, "class": {"id": '${class7Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)



enrollment='{"student": {"id": '${student8Id}'}, "class": {"id": '${class8Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student9Id}'}, "class": {"id": '${class9Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student10Id}'}, "class": {"id": '${class10Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student11Id}'}, "class": {"id": '${class1Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student12Id}'}, "class": {"id": '${class2Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student13Id}'}, "class": {"id": '${class3Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student14Id}'}, "class": {"id": '${class4Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student15Id}'}, "class": {"id": '${class5Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student16Id}'}, "class": {"id": '${class6Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student17Id}'}, "class": {"id": '${class7Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


enrollment='{"student": {"id": '${student18Id}'}, "class": {"id": '${class8Id}'}}'
request="curl -X POST localhost:8080/studentEnrollments -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${enrollment})'"
resp="$(eval ${request})"
error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
if [ -n "$error" ]; then
	clear
	echo ${error}
	echo ${request}
	exit
fi
enrollment1Id=$(echo ${resp} | cut -d ',' -f 1)
enrollment1Id=$(echo ${enrollment1Id} | cut -d ':' -f 2)


### Recados ###

# Recado para Segmentos
for ((i=segment1Id; i<=segment2Id; i++)); do
   	text=$(bash get_random_text.sh)
	message='{"title": "Racado para o segmento '${i}'", "text": "'${text}'", "registerType": {"id": 1}, "targetId": '${i}', "studentId": 0,
 "institutionId": '${institutionId}'}'
	request="curl -X POST localhost:8080/registers/segment -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${message})'"
	resp="$(eval ${request})"
	error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
	if [ -n "$error" ]; then
		clear
		echo ${error}
		echo ${request}
		exit
	fi
done

# Recado para Séries
for ((i=grade1Id; i<=grade6Id; i++)); do
   	text=$(bash get_random_text.sh)
	message='{"title": "Racado para a série '${i}'", "text": "'${text}'", "registerType": {"id": 1}, "targetId": '${i}', "studentId": 0, "institutionId": '${institutionId}'}'
	request="curl -X POST localhost:8080/registers/grade -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${message})'"
	resp="$(eval ${request})"
	error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
	if [ -n "$error" ]; then
		clear
		echo ${error}
		echo ${request}
		exit
	fi
done

# Recado para turmas
for ((i=class1Id; i<=class10Id; i++)); do
   	text=$(bash get_random_text.sh)
	message='{"title": "Racado para Turma '${i}'", "text": "'${text}'", "registerType": {"id": 1}, "targetId": '${i}', "studentId": 0, "institutionId": '${institutionId}'}'
	request="curl -X POST localhost:8080/registers/class -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${message})'"
	resp="$(eval ${request})"
	error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
	if [ -n "$error" ]; then
		clear
		echo ${error}
		echo ${request}
		exit
	fi
done

# Recado dos pais 

# Cada pai vai mandar 10 mensagens
for ((i=0; i<10; i++)); do
	for ((j=user1Id; j<=user10Id; j++)); do
		text=$(bash get_random_text.sh)
		title=$(echo ${text} | cut -d ',' -f 1 | cut -d '.' -f 1)
		echo ${title}
		message='{"title": "'${title}'", "text": "'${text}'", "registerType": {"id": 1}, "targetId": 1, "senderId": '${j}', "institutionId": '${institutionId}'}'
		request="curl -X POST localhost:8080/registers -H '"${auth}"' -H 'Content-Type: application/json' -d '$(echo ${message})'"
		resp="$(eval ${request})"
		error=$(echo ${resp} | awk -F 'error":' '{print $2}' | cut -d '}' -f 1)
		if [ -n "$error" ]; then
			clear
			echo ${error}
			echo ${request}
			exit
		fi
	done;
done 


echo "FIM DE ROTINA"