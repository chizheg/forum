// @jest-environment jsdom
const fs = require('fs');
const path = require('path');

describe('Forum Frontend Tests', () => {
    beforeEach(() => {
        document.body.innerHTML = fs.readFileSync(
            path.resolve(__dirname, './index.html'),
            'utf8'
        );
        localStorage.clear();
        window.WebSocket = jest.fn();
    });

    test('initial UI state shows login/register buttons', () => {
        expect(document.getElementById('auth-buttons').style.display).not.toBe('none');
        expect(document.getElementById('user-info').style.display).toBe('none');
    });

    test('login button shows login modal', () => {
        document.getElementById('login-btn').click();
        expect(document.getElementById('login-modal').classList.contains('hidden')).toBeFalsy();
    });

    test('register button shows register modal', () => {
        document.getElementById('register-btn').click();
        expect(document.getElementById('register-modal').classList.contains('hidden')).toBeFalsy();
    });

    test('modal close buttons hide modals', () => {
        document.getElementById('login-btn').click();
        document.querySelector('.modal-close').click();
        expect(document.getElementById('login-modal').classList.contains('hidden')).toBeTruthy();
    });

    test('successful login updates UI and connects WebSocket', async () => {
        global.fetch = jest.fn(() =>
            Promise.resolve({
                ok: true,
                json: () => Promise.resolve({ token: 'test-token' })
            })
        );

        const username = 'testuser';
        const password = 'testpass';

        document.getElementById('login-username').value = username;
        document.getElementById('login-password').value = password;
        
        await document.getElementById('login-submit').click();

        expect(localStorage.getItem('token')).toBe('test-token');
        expect(localStorage.getItem('username')).toBe(username);
        expect(document.getElementById('auth-buttons').style.display).toBe('none');
        expect(document.getElementById('user-info').style.display).toBe('block');
        expect(window.WebSocket).toHaveBeenCalled();
    });

    test('logout clears storage and updates UI', () => {
        localStorage.setItem('token', 'test-token');
        localStorage.setItem('username', 'testuser');
        
        document.getElementById('logout-btn').click();

        expect(localStorage.getItem('token')).toBeNull();
        expect(localStorage.getItem('username')).toBeNull();
        expect(document.getElementById('auth-buttons').style.display).toBe('block');
        expect(document.getElementById('user-info').style.display).toBe('none');
    });

    test('message sending requires WebSocket connection', () => {
        const ws = { send: jest.fn() };
        window.ws = ws;

        const messageText = document.getElementById('message-text');
        messageText.value = 'Test message';

        document.getElementById('send-btn').click();

        expect(ws.send).toHaveBeenCalledWith(JSON.stringify({
            content: 'Test message'
        }));
        expect(messageText.value).toBe('');
    });
}); 