'use strict';

document.getElementById('addMore').addEventListener('click', () => {
    let newInput = document.createElement('input');
    newInput.name = 'array[]';
    newInput.type = 'text';
    newInput.placeholder = 'Enter Host Name';
    document.getElementById('arrs').appendChild(newInput);
});
const form = document.getElementById('myForm');
form.addEventListener('submit', (e) => {
    e.preventDefault();
    let formValues = [];
    let inputElements = document.querySelectorAll('#arrs > input[name="array[]"]');
    inputElements.forEach((elm) => {
        if (!/^\s*$/gi.test(elm.value)) {
            formValues.push(elm.value);
        };
    });
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