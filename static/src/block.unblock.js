'use strict';
const result = document.getElementById('result');
document.getElementById('mybtn').addEventListener('click', () => {
    result.innerHTML = 'Please Weait...';
    fetch(url).then(res => res.json()).then(data => {
        if (data.msg === 'success') {
            result.innerHTML = 'success';
        } else {
            result.innerHTML = JSON.stringify(data, null, '   ');
        };
    }).catch(err => {
        console.error(err);
    });
});