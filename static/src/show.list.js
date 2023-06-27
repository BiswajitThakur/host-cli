'use strict';

document.getElementById('searchBtn').addEventListener('click', () => {
    let val = document.getElementById('searchInput').value.replace(/[\\\/\|\+\.\-]/gi, (e) => {
        return '\\' + e;
    });
    document.querySelectorAll('.aa1 > .bb1').forEach(element => {
        if (new RegExp(val, 'gi').test(element.textContent)) {
            element.style.display = 'block';
        } else {
            element.style.display = 'none';
        };
    });
});
let aa1 = document.getElementById('aa1');
const ld = document.getElementById('ld');
fetch(fetchUrl).then(res => res.json()).then(data => {
    for (let i = 0; i < data.length; i++) {
        let elm = document.createElement('div');
        elm.setAttribute('class', 'bb1');
        elm.textContent = data[i];
        aa1.appendChild(elm);
    };
    ld.style.display = 'none';
}).catch(err => {
    ld.innerHTML = err;
    console.error(err);
});