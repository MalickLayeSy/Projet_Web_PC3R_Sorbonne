let registerModal;
let loginModal;
let loginBtn;
let registerBtn;
let registerForm;
let loginForm;
let idSession = "";
let pariEnCoursListe;
let matchDisponibleListe;
let matchDisponibleButton;

function majSiIdSession() {
    if (idSession !== "" && idSession !== null) {
        sessionStorage.setItem('idSession', idSession);
        document.getElementById("loginBtn").style.display = "none"
        document.getElementById("registerBtn").style.display = "none"
        document.getElementById("disconnectBtn").style.display = "block"
        document.getElementById("montantCompte").style.display = "block"
        refreshMatchComing()
        refreshActiveBet()
        refreshBetHistory()
        changeUserMoney()
    }
}

window.onload = function () {
    registerModal = document.getElementById("registerModal");
    loginModal = document.getElementById("loginModal");
    loginBtn = document.getElementById("loginBtn");
    registerBtn = document.getElementById("registerBtn");
    registerForm = document.getElementById("registerForm")
    loginForm = document.getElementById("loginForm")
    pariEnCoursListe = document.getElementById("pariEnCoursListe")
    matchDisponibleListe = document.getElementById("matchDisponibleListe")
    matchDisponibleButton = document.getElementById("matchDisponibleButton")

    idSession = sessionStorage.getItem('idSession');
    console.log("ID SESSION " + idSession)
    majSiIdSession();

    loginForm.addEventListener('submit', function (event) {
        event.preventDefault();
        window.location.hash = "#home"

        let formData = new FormData(event.target);
        fetch('/connexion', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                'Accept': 'application/json'
            },
            body: new URLSearchParams(formData)
        })
            .then(function (response) {
                return response.json();
            })
            .then(function (jsonData) {
                window.alert(jsonData["message"])
                if (jsonData["code"] === "200") {
                    loginModal.style.display = "none";
                    clearLoginForm()
                    idSession = jsonData["idSession"]
                    majSiIdSession();
                }
            });
    })

    document.getElementById("accueil").onclick = function () {
        window.location.hash = "#home"

        document.getElementById("main").style.display = 'block'
        document.getElementById("resultSearch").style.display = 'none'
    }

    document.getElementById("searchBtn").onclick = function () {
        let params = new URLSearchParams()
        params.append("idSession", idSession)
        params.append("req", document.getElementById("site-search").value)
        fetch("/match?" + params.toString())
            .then(res => res.json())
            .then(function (jsonData) {
                if (jsonData["code"] === "200") {
                    document.getElementById("resultSearch").style.display = "block"

                    document.getElementById("main").style.display = 'none'

                    const resultList = document.getElementById('resultSearchList');

                    window.location.hash = "#search"

                    resultList.innerHTML = ""


                    let result = jsonData["result"]

                    for (let i = 0; i < result.length; i++) {
                        const form = document.createElement('form')
                        form.name = result[i]["id"]

                        const submitButton = document.createElement('button')

                        submitButton.type = "submit"

                        submitButton.appendChild(document.createTextNode("Soumettre Pari"))

                        const li = document.createElement('li');

                        const ul = document.createElement('ul');

                        const sport = document.createTextNode("Sport : " + result[i]["sport"]);
                        const sportLi = document.createElement('li');
                        sportLi.appendChild(sport)

                        const league = document.createTextNode("League : " + result[i]["league"]);
                        const leagueLi = document.createElement('li');
                        leagueLi.appendChild(league)

                        const date = document.createTextNode('Date : ' + result[i]["date"]);
                        const dateLi = document.createElement('li');
                        dateLi.appendChild(date)

                        const equipe = document.createTextNode(result[i]["equipeA"] + " vs " + result[i]["equipeB"]);
                        const equipeLi = document.createElement('li');
                        equipeLi.appendChild(equipe)

                        const cote = document.createTextNode("Cote : " + result[i]["cote"]);
                        const coteLi = document.createElement('li');
                        coteLi.appendChild(cote)

                        let montant;
                        if (result[i]["vainqueur"] === "") {
                           montant = document.createElement('input');
                            const montantLi = document.createElement('li');
                            const montantTxt = document.createTextNode("Montant ");
                            montantLi.append(montantTxt, montant)
                            montant.type = 'number'
                            montant.value = 0
                            montant.min = 0

                            const vainqueurLi = document.createElement('li');

                            const equipeARadio = document.createElement('input');
                            equipeARadio.type = 'radio'
                            equipeARadio.value = result[i]["equipeA"]
                            equipeARadio.name = 'vainqueur' + result[i]["id"]
                            equipeARadio.id = 'equipeA' + result[i]["id"]

                            const equipeALabel = document.createElement('label');
                            equipeALabel.htmlFor = 'equipeA' + result[i]["id"]

                            const equipeALabelText = document.createTextNode(result[i]["equipeA"]);

                            equipeALabel.appendChild(equipeALabelText)

                            const equipeBRadio = document.createElement('input');
                            equipeBRadio.type = 'radio'
                            equipeBRadio.value = result[i]["equipeB"]
                            equipeBRadio.name = 'vainqueur' + result[i]["id"]
                            equipeBRadio.id = 'equipeB' + result[i]["id"]

                            const equipeBLabel = document.createElement('label');
                            equipeBLabel.htmlFor = 'equipeB' + result[i]["id"]

                            const equipeBLabelText = document.createTextNode(result[i]["equipeB"]);

                            equipeBLabel.appendChild(equipeBLabelText)

                            vainqueurLi.appendChild(equipeARadio)
                            vainqueurLi.appendChild(equipeALabel)
                            vainqueurLi.appendChild(equipeBRadio)
                            vainqueurLi.appendChild(equipeBLabel)

                            ul.appendChild(sportLi)
                            ul.appendChild(leagueLi)
                            ul.appendChild(dateLi)
                            ul.appendChild(equipeLi)
                            ul.appendChild(coteLi)
                            ul.appendChild(vainqueurLi)
                            ul.appendChild(montantLi)

                            form.appendChild(ul)
                            form.appendChild(submitButton)
                        } else {
                            const vainqueurLi = document.createElement('li');
                            const vainqueurLabelText = document.createTextNode("Vainqueur : " + result[i]["vainqueur"]);

                            vainqueurLi.appendChild(vainqueurLabelText)

                            ul.appendChild(sportLi)
                            ul.appendChild(leagueLi)
                            ul.appendChild(dateLi)
                            ul.appendChild(equipeLi)
                            ul.appendChild(coteLi)
                            ul.appendChild(vainqueurLi)

                            form.appendChild(ul)
                            // form.appendChild(submitButton)
                        }



                        li.appendChild(form)

                        resultList.append(li)

                        submitButton.onclick = function (event) {
                            event.preventDefault()
                            let idMatch = event.target.form.name
                            let vainqueur = document.querySelector('input[name="vainqueur' + idMatch + '"]:checked').value;
                            if (montant.value !== 0 && vainqueur !== "") {
                                params = new URLSearchParams()
/*                                console.log("Montant : " + montant.value)
                                console.log("Cote : " + result[i]["cote"])*/
                                params.append("idSession", idSession)
                                params.append("idMatch", idMatch)
                                params.append("equipeGagnante", vainqueur)
                                params.append("cote", result[i]["cote"])
                                params.append("montant", montant.value)

                                fetch("/bet", {
                                    method: 'POST',
                                    headers: {
                                        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                                        'Accept': 'application/json'
                                    },
                                    body: params
                                })
                                    .then(function (response) {
                                        // console.log(response)
                                        return response.json()
                                    })
                                    .then(function (jsonData) {
                                        window.alert(jsonData["message"])
                                    })
                            }
                        }
                    }

                }
            })

    }

    document.getElementById("disconnectBtn").onclick = function () {
        console.log("Déconnexion")
        let params = new URLSearchParams()
        params.append("idSession", idSession)
        fetch("/connexion?" + params.toString(), {
            method: 'DELETE',
        })
            .then(res => res.json())
            .then(function (jsonData) {

                if (jsonData["code"] === "200") {
                    alert(jsonData["message"])
                    document.getElementById("loginBtn").style.display = "block"
                    document.getElementById("registerBtn").style.display = "block"
                    document.getElementById("disconnectBtn").style.display = "none"
                    document.getElementById("montantCompte").style.display = "none"
                    clearChamp('pariEnCoursListe')
                    clearChamp('matchDisponibleListe')
                    clearChamp('pariFinisListe')
                }
            })
    }

    registerForm.addEventListener('submit', function (event) {
        event.preventDefault();
        let formData = new FormData(event.target);

        fetch('/user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                'Accept': 'application/json'
            },
            body: new URLSearchParams(formData)
        })
            .then(function (response) {
                return response.json();
            })
            .then(function (jsonData) {
                // console.log(jsonData);
                window.alert(jsonData["message"])
            });
    })


    collapse()


    loginBtn.onclick = function () {
        loginModal.style.display = "block";
    }

    registerBtn.onclick = function () {
        registerModal.style.display = "block";
    }
}

function clearRegisterForm() {
    document.getElementById("loginRegister").value = ''
    document.getElementById("mailRegister").value = ''
    document.getElementById("passwordRegister").value = ''
}

function clearLoginForm() {
    document.getElementById("loginConnect").value = ''
    document.getElementById("passwordConnect").value = ''
}

function refreshMatchComing() {

    let params = new URLSearchParams()
    params.append("idSession", idSession)
    // console.log(params.toString())
    fetch("/match?" + params.toString())
        .then(function (response) {
            return response.json();
        })
        .then(function (jsonData) {
            if (jsonData["code"] === "200") {
                clearMatchAVenir()
                let result = jsonData["result"];
                for (let i = 0; i < result.length; i++) {
                    const form = document.createElement('form')
                    form.name = result[i]["id"]

                    const submitButton = document.createElement('button')

                    submitButton.type = "submit"

                    submitButton.appendChild(document.createTextNode("Soumettre Pari"))

                    const li = document.createElement('li');

                    const ul = document.createElement('ul');

                    const sport = document.createTextNode("Sport : " + result[i]["sport"]);
                    const sportLi = document.createElement('li');
                    sportLi.appendChild(sport)

                    const league = document.createTextNode("League : " + result[i]["league"]);
                    const leagueLi = document.createElement('li');
                    leagueLi.appendChild(league)

                    const date = document.createTextNode('Date : ' + result[i]["date"]);
                    const dateLi = document.createElement('li');
                    dateLi.appendChild(date)

                    const equipe = document.createTextNode(result[i]["equipeA"] + " vs " + result[i]["equipeB"]);
                    const equipeLi = document.createElement('li');
                    equipeLi.appendChild(equipe)

                    let coteA = result[i]["cote"]
                    let coteB = 100 / coteA

                    const cotePourA = document.createTextNode("Cote pour l'équipe " + result[i]["equipeA"] + " : " + coteA);
                    const coteLiA = document.createElement('li');
                    coteLiA.appendChild(cotePourA)

                    const cotePourB = document.createTextNode("Cote pour l'équipe " + result[i]["equipeB"] + " : " + coteB);
                    const coteLiB = document.createElement('li');
                    coteLiB.appendChild(cotePourB)

                    const montant = document.createElement('input');
                    const montantLi = document.createElement('li');
                    const montantTxt = document.createTextNode("Montant ");
                    montantLi.append(montantTxt, montant)
                    montant.type = 'number'
                    montant.value = 0
                    montant.min = 0

                    const vainqueurLi = document.createElement('li');

                    const equipeARadio = document.createElement('input');
                    equipeARadio.type = 'radio'
                    equipeARadio.value = result[i]["equipeA"]
                    equipeARadio.name = 'vainqueur' + result[i]["id"]
                    equipeARadio.id = 'equipeA' + result[i]["id"]

                    const equipeALabel = document.createElement('label');
                    equipeALabel.htmlFor = 'equipeA' + result[i]["id"]

                    const equipeALabelText = document.createTextNode(result[i]["equipeA"]);

                    equipeALabel.appendChild(equipeALabelText)

                    const equipeBRadio = document.createElement('input');
                    equipeBRadio.type = 'radio'
                    equipeBRadio.value = result[i]["equipeB"]
                    equipeBRadio.name = 'vainqueur' + result[i]["id"]
                    equipeBRadio.id = 'equipeB' + result[i]["id"]

                    const equipeBLabel = document.createElement('label');
                    equipeBLabel.htmlFor = 'equipeB' + result[i]["id"]

                    const equipeBLabelText = document.createTextNode(result[i]["equipeB"]);

                    equipeBLabel.appendChild(equipeBLabelText)

                    vainqueurLi.appendChild(equipeARadio)
                    vainqueurLi.appendChild(equipeALabel)
                    vainqueurLi.appendChild(equipeBRadio)
                    vainqueurLi.appendChild(equipeBLabel)

                    ul.appendChild(sportLi)
                    ul.appendChild(leagueLi)
                    ul.appendChild(dateLi)
                    ul.appendChild(equipeLi)
                    ul.appendChild(coteLiA)
                    ul.appendChild(coteLiB)
                    ul.appendChild(vainqueurLi)
                    ul.appendChild(montantLi)

                    form.appendChild(ul)
                    form.appendChild(submitButton)

                    li.appendChild(form)

                    matchDisponibleListe.append(li)

                    submitButton.onclick = function (event) {
                        event.preventDefault()
                        let idMatch = event.target.form.name
                        let vainqueur = document.querySelector('input[name="vainqueur' + idMatch + '"]:checked').value;
                        if (montant.value !== 0 && vainqueur !== "") {
                            params = new URLSearchParams()
                            // console.log("Montant : " + montant.value)
                            // console.log("Cote : " + result[i]["cote"])
                            params.append("idSession", idSession)
                            params.append("idMatch", idMatch)
                            params.append("equipeGagnante", vainqueur)
                            params.append("cote", result[i]["cote"])
                            params.append("montant", montant.value)

                            fetch("/bet", {
                                method: 'POST',
                                headers: {
                                    'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                                    'Accept': 'application/json'
                                },
                                body: params
                            })
                                .then(function (response) {
                                    // console.log(response)
                                    return response.json()
                                })
                                .then(function (jsonData) {
                                    window.alert(jsonData["message"])
                                })
                        }
                    }
                }
            } else {
                alert(jsonData["message"])
            }
        });
}

function refreshActiveBet() {
    getBet("coming", "pariEnCoursListe")

}

function refreshBetHistory() {
    getBet("", "pariFinisListe")
}

function getBet(statut, champ) {
    let params = new URLSearchParams()
    params.append("idSession", idSession)
    params.append("statutParis", statut)
    fetch("/bet?" + params.toString())
        .then(function (response) {
            return response.json()
        })
        .then(function (jsonData) {
            console.log(jsonData["message"])
            if (jsonData["code"] === "200") {
                clearChamp(champ)
                let result = jsonData["result"]
                for (let i = 0; i < result.length; i++) {

                    const li = document.createElement('li')

                    const ul = document.createElement('ul')

                    Object.keys(result[i]).forEach(function (key) {
                        const value = result[i][key];

                        let liUl = document.createElement('li')

                        let texte = document.createTextNode(key + " : " + value)

                        liUl.appendChild(texte)

                        ul.appendChild(liUl)
                    })

                    li.appendChild(ul)

                    document.getElementById(champ).append(li)
                }
            } else {
                alert(jsonData["message"])
            }
        })
}

window.onclick = function (event) {
    if (event.target === registerModal) {
        registerModal.style.display = "none";
        clearRegisterForm()
    } else if (event.target === loginModal) {
        loginModal.style.display = "none";
        clearLoginForm()
    }
}

function clearMatchAVenir() {
    document.getElementById("matchDisponibleListe").innerHTML = "";
}

function clearChamp(champ) {
    document.getElementById(champ).innerHTML = "";
}

function collapse() {
    const coll = document.getElementsByClassName("collapsible");
    let i;

    for (i = 0; i < coll.length; i++) {
        coll[i].addEventListener("click", function () {
            this.classList.toggle("active");
            const content = this.nextElementSibling;
            if (content.style.display === "block") {
                content.style.display = "none";
            } else {
                content.style.display = "block";
            }
        });
    }
}

function changeUserMoney() {
    let params = new URLSearchParams()
    params.append("idSession", idSession)
    fetch("/user?" + params.toString())
        .then(res => res.json()
        )
        .then(function (jsonData) {
            let cagnotte = Number.parseFloat(jsonData["cagnotte"]).toFixed(2)
            document.getElementById("montantCompte").innerText = "Montant compte : " + cagnotte
        })
}

window.onpopstate = function () {
    switch ( location.hash ) {
        case '#home':
            backFromHome()
            break
        case '#search':
            backFromSearch()
            break
        default:
        // defaultBackAnimation()
    }
}

function backFromHome() {
    document.getElementById("main").style.display = "block"
    document.getElementById("resultSearch").style.display = "none"
}

function backFromSearch() {
    document.getElementById("main").style.display = "none"
    document.getElementById("resultSearch").style.display = "block"
}
