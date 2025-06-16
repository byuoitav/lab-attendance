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

    // after 30 seconds of inactivity, show the screensaver
    let inactivityTimeout;
    function resetInactivityTimeout() {
        window.hideScreensaver();
        clearTimeout(inactivityTimeout);
        inactivityTimeout = setTimeout(() => {
            window.showScreensaver();
        }, 30000); // 30 seconds
    }

    // Reset inactivity timeout on user interaction
    document.addEventListener('mousemove', resetInactivityTimeout);
    document.addEventListener('keydown', resetInactivityTimeout);
    resetInactivityTimeout(); // Initialize the timeout

    
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
    // Remove screensaver if it's showing
    window.hideScreensaver();
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

window.showScreensaver = function() {
    // Remove any existing screensaver first
    window.hideScreensaver();

    const screensaver = document.createElement('div');
    screensaver.id = 'screensaver';
    screensaver.className = 'screensaver';

    const time = document.createElement('h1');
    time.id = 'screensaver-time';
    time.textContent = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    screensaver.appendChild(time);

    const message = document.createElement('p');
    message.textContent = 'Sign or Tap into the Lab.';
    screensaver.appendChild(message);

    document.body.appendChild(screensaver);

    // Start interval to update time every second
    window.screensaverTimeInterval = setInterval(() => {
        const timeElem = document.getElementById('screensaver-time');
        if (timeElem) {
            timeElem.textContent = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        }
    }, 1000);
};

// Clear the interval when hiding the screensaver
window.hideScreensaver = function() {
    const screensaver = document.getElementById('screensaver');
    if (screensaver) {
        screensaver.remove();
    }
    if (window.screensaverTimeInterval) {
        clearInterval(window.screensaverTimeInterval);
        window.screensaverTimeInterval = null;
    }
};

window.hideScreensaver = function() {
    const screensaver = document.getElementById('screensaver');
    if (screensaver) {
        screensaver.remove();
    }
};