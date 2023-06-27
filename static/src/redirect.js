'use strict';
document.getElementById('addMore').addEventListener('click', () => {
    let newInput = document.createElement('input');
    newInput.name = 'array[]';
    newInput.type = 'text';
    newInput.required = 'required';
    newInput.setAttribute('class', 'inp_main');
    newInput.placeholder = 'Enter Host Name';
    document.getElementById('arrs').appendChild(newInput);
    let newInput1 = document.createElement('input');
    newInput1.name = 'array[]';
    newInput1.type = 'text';
    newInput1.required = 'required';
    newInput1.setAttribute('class', 'inp_main');
    newInput1.placeholder = 'Enter Host Name';
    document.getElementById('arrs').appendChild(newInput1);
});
const form = document.getElementById('myForm');
form.addEventListener('submit', (e) => {
    e.preventDefault();
    let formValues = [];
    let inputElements = document.querySelectorAll('#arrs > input[name="array[]"]');
    for (let i = 0; i < inputElements.length; i = i + 2) {
        if (/^\S*$/gi.test(inputElements[i].value) && /^\S*$/gi.test(inputElements[i + 1].value)) {
            formValues.push([inputElements[i].value, inputElements[i + 1].value]);
        };
    };
    if (formValues.length === 0) {
        return;
    };
    let resultElm = document.getElementById('result');
    resultElm.innerHTML = 'please wait....';
    fetch(form.action, {
        method: 'POST',
        body: JSON.stringify(formValues),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(response => response.json()).then((data) => {
        if (data.msg === 'success') {
            resultElm.innerHTML = 'success';
            setTimeout(() => {
                resultElm.innerHTML = '';
                inputElements.forEach((elm) => {
                    elm.value = '';
                });
            }, 4000);
        } else {
            resultElm.innerHTML = 'failed';
            setTimeout(() => {
                resultElm.innerHTML = '';
            }, 4000);
        };
    }).catch(error => console.error(error));
});