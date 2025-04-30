async function addAchievement() {
    const form = document.querySelector("#addAchievement");
    const data = new FormData(form);

    var object = {};
    data.forEach(function(value, key){
        object[key] = value;
    });
    var json = JSON.stringify(object);
    console.log(json)

    try {
        const token = getCookie('receivedToken');

        if (!token) {
            throw new Error('Token non trouvé dans les cookies');
        }

        const response = await fetch('/admin/achievements/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `${token}`,
            },
            body: json,
        });

        if (!response.ok) {
            throw new Error('Échec de l\'ajout de l\'achievement');
        }

    } catch (error) {
        console.error('Erreur lors de l\'ajout de l\'achievement :', error);
    }
}

async function deleteAchievement(id){
    try {
        const token = getCookie('receivedToken');

        if (!token) {
            throw new Error('Token non trouvé dans les cookies');
        }

        const response = await fetch('/admin/achievements/'+ id, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `${token}`,
            },
        });

        if (!response.ok) {
            throw new Error('Échec de suppression de l\'achievement');
        }

    } catch (error) {
        console.error('Erreur lors de la suppression de l\'achievement :', error);
    }
}

async function updateAchievement(){
    const form = document.querySelector("#addAchievement");
    const data = new FormData(form);

    id = data.get('id');
    var object = {};
    data.forEach(function(value, key){
        if (key !== 'id') {
            object[key] = value;
        }
    });
    var json = JSON.stringify(object);
    console.log(json)


    try {
        const token = getCookie('receivedToken');

        if (!token) {
            throw new Error('Token non trouvé dans les cookies');
        }

        const response = await fetch('/admin/achievements/'+id, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `${token}`,
            },
            body: json,
        });

        if (!response.ok) {
            throw new Error('Échec de l\'ajout de l\'achievement');
        }

    } catch (error) {
        console.error('Erreur lors de l\'ajout de l\'achievement :', error);
    }
}

function getCookie(name) {
    const cookies = document.cookie.split(';');
    for (let cookie of cookies) {
        if (cookie.startsWith(name + '=')) {
            return cookie.substring(name.length + 1);
        }
    }

    return null;
}