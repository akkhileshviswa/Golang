let subscribeBtn = document.getElementById("subscribe");
let form = document.getElementById("form");
let brokerBtn = document.getElementById("broker");

function validate() {
    var firstName = validateFirstName();
    var lastName = validateLastName();
    var email = validateEmail();

    if (firstName && lastName && email) {
        return true;
    }
    return false;

    function validateFirstName() {
        var firstName = document.getElementById("firstName").value;
        var firstNameErr = document.getElementById("firstNameErr");
        var error = "";
        if (firstName == "") {
            error = "Please enter the first name";
        }
        firstNameErr.innerHTML = error;
        return error == "" ? true : false;
    }

    function validateLastName() {
        var lastName = document.getElementById("lastName").value;
        var lastNameErr = document.getElementById("lastNameErr");
        var error = "";
        if (lastName == "") {
            error = "Please enter the last name";
        }
        lastNameErr.innerHTML = error;
        return error == "" ? true : false;
    }

    function validateEmail() 
    {
        var email = document.getElementById("email").value;
        var emailErr = document.getElementById("emailErr");
        var error = "";
        var atposition = email.indexOf("@");  
        var dotposition = email.lastIndexOf(".");  
        if (email == "") {
            error = "Please enter the email";
        } else if (atposition < 1 || dotposition < atposition+2 || dotposition+2 >= email.length) {
            error = "Please enter the valid email"
        }
        emailErr.innerHTML = error;
        return error == "" ? true : false;
    }
}

subscribeBtn.addEventListener("click", function(e) {
    e.preventDefault();
    if(!validate()) {
        return false;
    }

    const payload = {
        action: "subscribe",
        submit: {
            firstName: document.getElementById("firstName").value,
            lastName: document.getElementById("lastName").value,
            email: document.getElementById("email").value,
        }
    }
    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    const body = {
        method: 'POST',
        body: JSON.stringify(payload),
        headers: headers,
    }

    fetch("http:\/\/localhost:8080/subscribe", body)
    .then((response) => response.json())
    .then((data) => {
        if (data.error) {
            form.innerHTML += `<br><strong>`+ data.message +`</strong>`;
        } else {
            form.innerHTML += `<br><strong>Thank you for signing up for the newsletter.<br>
                Kindly check your mail for more updates!</strong>`;
        }
    })
    .catch((error) => {
        form.innerHTML += "<br><br>Errortry: " + error;
    })
})
