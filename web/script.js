document.addEventListener('DOMContentLoaded', async () => {
    currentComponent = 'keypad';
    loadComponent(currentComponent);

    // make the services available globally
    window.apiService = await new APIService();
    window.eventService = new EventService();

    // Fetch lab name from config endpoint and set it
    const labName = await window.apiService.getLabName();
    if (labName) {
        document.querySelector('.lab-name').textContent = labName;
    }
});

async function loadComponent(name) {
    // call cleanup on the current component
    if (window.components?.[currentComponent]?.cleanup) {
        window.components[currentComponent].cleanup();
    }

    const htmlPath = `./${name}/${name}.html`;
    const jsPath = `./${name}/${name}.js`;
    const cssPath = `./${name}/${name}.css`;

    // load the css
    const oldStylesheet = document.getElementById('component-stylesheet');
    if (oldStylesheet) {
        oldStylesheet.remove();
    }

    const stylesheet = document.createElement('link');
    stylesheet.rel = 'stylesheet';
    stylesheet.href = cssPath;
    stylesheet.id = 'component-stylesheet';
    stylesheet.onload = () => {
        const module = window.components?.[name];
        if (module?.loadStyles) {
            module.loadStyles();
        }
    }

    document.head.appendChild(stylesheet);

    // load the html
    const componentContainer = document.querySelector('.content-right');
    componentContainer.classList.add('loading'); // hide before loading

    const response = await fetch(htmlPath);
    const html = await response.text();
    componentContainer.innerHTML = html;

    // load the js
    const oldScript = document.getElementById('component-script');
    if (oldScript) {
        oldScript.remove();
    }

    const script = document.createElement('script');
    script.src = jsPath;
    script.id = 'component-script';
    // call loadPage on the new component
    script.onload = () => {
        const module = window.components?.[name];
        if (module?.loadPage) {
            module.loadPage();
            currentComponent = name;
        }
        componentContainer.classList.remove('loading'); // finally show it
    };

    document.body.appendChild(script);
}

// Show a popup for welcome or error messages
window.showPopup = function(data) {
    console.log('showPopup called with data:', data);
    console.log('data.key:', data?.key);
    // Remove any existing popup
    const oldPopup = document.getElementById('popup');
    if (oldPopup) oldPopup.remove();

    const popup = document.createElement('div');
    popup.id = 'popup';
    popup.className = 'popup';

    let messages = [];
    switch (data.key) {
        case 'login':
            messages = [
                `Welcome, ${data.data.FirstName}!`,
                `Successfully logged in ${data.data.Name}`
            ];
            break;
        case 'login-error':
            messages = [
                'Login failed:',
                data.data
            ];
            break;
        case 'card-read-error':
            messages = [
                'Unable to read card, please try your ID Card again.'
            ];
            break;
        default:
            messages = ['Test'];
    }
    messages.forEach(msg => {
        const p = document.createElement('p');
        p.textContent = msg;
        popup.appendChild(p);
    });

    document.body.appendChild(popup);
    setTimeout(() => {
        popup.remove();
    }, 3500);
};
