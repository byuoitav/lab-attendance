window.components = window.components || {};
window.components.keypad = {
    loadPage: function () {
        this.loadKeypad();
    },

    cleanup: function () {
        this.idEntry = null;
    },

    loadKeypad: function () {
        this.idEntry = document.querySelector('.id-entry');
        this.backspaceBtn = document.querySelector('.del');
        this.enterBtn = document.querySelector('.enter');
        const keys = [
            { selector: '.one', value: '1' },
            { selector: '.two', value: '2' },
            { selector: '.three', value: '3' },
            { selector: '.four', value: '4' },
            { selector: '.five', value: '5' },
            { selector: '.six', value: '6' },
            { selector: '.seven', value: '7' },
            { selector: '.eight', value: '8' },
            { selector: '.nine', value: '9' },
            { selector: '.zero', value: '0' }
        ];
        keys.forEach(key => {
            const btn = document.querySelector(key.selector);
            if (btn) btn.addEventListener('click', () => this.appendIdEntry(key.value));
        });
        if (this.backspaceBtn) this.backspaceBtn.addEventListener('click', () => this.backspaceIdEntry());
        if (this.enterBtn) this.enterBtn.addEventListener('click', () => {
            const byuId = this.idEntry.textContent.replace(/-/g, ''); // Remove hyphens
            window.apiService.login(byuId);
            this.clearIdEntry()
        });
        this.updateButtonStates();
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
            this.updateButtonStates();
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
            this.updateButtonStates();
        } else {
            console.error('.id-entry not found');
        }
    },

    clearIdEntry: function () {
        if (this.idEntry) {
            this.idEntry.textContent = 'Enter BYU ID'; // Reset to placeholder text
            this.updateButtonStates();
        } else {
            console.error('.id-entry not found');
        }
    },

    updateButtonStates: function () {
        const entryLength = this.idEntry && this.idEntry.textContent && this.idEntry.textContent !== 'Enter BYU ID'
            ? this.idEntry.textContent.replace(/-/g, '').length
            : 0;
        if (this.backspaceBtn) {
            if (entryLength > 0) {
                this.backspaceBtn.classList.remove('unclickable');
            } else {
                this.backspaceBtn.classList.add('unclickable');
            }
        }
        if (this.enterBtn) {
            if (entryLength === 9) {
                this.enterBtn.classList.remove('unclickable');
            } else {
                this.enterBtn.classList.add('unclickable');
            }
        }
    }
}
