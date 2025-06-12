window.components = window.components || {};
window.components.keypad = {
    loadPage: function () {
        this.loadKeypad();
    },

    cleanup: function () {
        console.log('nothing to cleanup');
    },

    loadKeypad: function () {
        console.log('Keypad script loaded');
        this.idEntry = document.querySelector('.id-entry');  // <-- assign to this.idEntry

        const one = document.querySelector('.one');
        const two = document.querySelector('.two');
        const three = document.querySelector('.three');
        const four = document.querySelector('.four');
        const five = document.querySelector('.five');
        const six = document.querySelector('.six');
        const seven = document.querySelector('.seven');
        const eight = document.querySelector('.eight');
        const nine = document.querySelector('.nine');
        const zero = document.querySelector('.zero');
        const backspace = document.querySelector('.del');
        const enter = document.querySelector('.enter');

        if (one) one.addEventListener('click', () => this.appendIdEntry('1'));
        if (two) two.addEventListener('click', () => this.appendIdEntry('2'));
        if (three) three.addEventListener('click', () => this.appendIdEntry('3'));
        if (four) four.addEventListener('click', () => this.appendIdEntry('4'));
        if (five) five.addEventListener('click', () => this.appendIdEntry('5'));
        if (six) six.addEventListener('click', () => this.appendIdEntry('6'));
        if (seven) seven.addEventListener('click', () => this.appendIdEntry('7'));
        if (eight) eight.addEventListener('click', () => this.appendIdEntry('8'));
        if (nine) nine.addEventListener('click', () => this.appendIdEntry('9'));
        if (zero) zero.addEventListener('click', () => this.appendIdEntry('0'));

        if (backspace) backspace.addEventListener('click', () => this.backspaceIdEntry());
        if (enter) enter.addEventListener('click', () => this.clearIdEntry());
    },

    appendIdEntry: function (text) {
        if (this.idEntry) {
            if (this.idEntry.textContent === 'Enter BYU ID') {
                this.idEntry.textContent = ''; // Clear placeholder text
            }

            if (this.idEntry.textContent.length >= 11) {
                console.warn('ID entry is full, cannot append more characters.');
                return;
            }

            if (this.idEntry.textContent.length === 2 || this.idEntry.textContent.length === 6) {
                this.idEntry.textContent += '-'; // Add hyphen after 3rd and 7th characters
            }

            this.idEntry.textContent = this.idEntry.textContent + text;
        } else {
            console.error('.id-entry not found');
        }
    },

    backspaceIdEntry: function () {
        if (this.idEntry) {
            if (this.idEntry.textContent === 'Enter BYU ID') {
                return; // Do nothing if placeholder text is present
            }
            this.idEntry.textContent = this.idEntry.textContent.slice(0, -1);
            // if last character is a hyphen, remove it
            if (this.idEntry.textContent.endsWith('-')) {
                this.idEntry.textContent = this.idEntry.textContent.slice(0, -1);
            }
            if (this.idEntry.textContent === '') {
                this.idEntry.textContent = 'Enter BYU ID'; // Clear placeholder text
            }
        } else {
            console.error('.id-entry not found');
        }
    },
// clairvoyant means 
    clearIdEntry: function () {
        if (this.idEntry) {
            this.idEntry.textContent = 'Enter BYU ID'; // Reset to placeholder text
        } else {
            console.error('.id-entry not found');
        }
    }
}

