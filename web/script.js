document.addEventListener('DOMContentLoaded', function() {
    const chatMessages = document.getElementById('chat-messages');
    const messageInput = document.getElementById('message-input');
    const sendButton = document.getElementById('send-button');

    const messages = [
        { role: 'system', content: "You are a helpful assistant that helps users manage your Cal.com calendar. You can help them book, list, and cancel events." }
    ];

    // Function to format bot messages
    function formatBotMessage(message) {
        // Convert numbered lists to <ol>
        if (/\d+\.\s/.test(message)) {
            const lines = message.split(/\n|<br\s*\/?>/i);
            let inList = false;
            let formatted = '';
            lines.forEach(line => {
                const match = line.match(/^(\d+)\.\s+(.*)/);
                if (match) {
                    if (!inList) {
                        formatted += '<ol>';
                        inList = true;
                    }
                    formatted += `<li>${match[2]}</li>`;
                } else {
                    if (inList) {
                        formatted += '</ol>';
                        inList = false;
                    }
                    formatted += line + '<br>';
                }
            });
            if (inList) formatted += '</ol>';
            return formatted;
        }
        // Convert dash lists to <ul>
        if (/^-\s+/m.test(message)) {
            const lines = message.split(/\n|<br\s*\/?>/i);
            let inList = false;
            let formatted = '';
            lines.forEach(line => {
                const match = line.match(/^[-*]\s+(.*)/);
                if (match) {
                    if (!inList) {
                        formatted += '<ul>';
                        inList = true;
                    }
                    formatted += `<li>${match[1]}</li>`;
                } else {
                    if (inList) {
                        formatted += '</ul>';
                        inList = false;
                    }
                    formatted += line + '<br>';
                }
            });
            if (inList) formatted += '</ul>';
            return formatted;
        }
        return message.replace(/\n/g, '<br>');
    }

    // Function to add a message to the chat
    function addMessage(message, isUser = false) {
        const messageElement = document.createElement('div');
        messageElement.classList.add('message');
        messageElement.classList.add(isUser ? 'user' : 'bot');

        const messageContent = document.createElement('div');
        messageContent.classList.add('message-content');
        messageContent.innerHTML = isUser ? message : formatBotMessage(message);

        messageElement.appendChild(messageContent);
        chatMessages.appendChild(messageElement);
        
        // Scroll to bottom of chat
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }

    // Function to send a message to the chatbot
    async function sendMessage(message) {
        if (!message.trim()) return;

        // Add user message to chat
        addMessage(message, true);
        messages.push({ role: 'user', content: message });

        // Clear input field
        messageInput.value = '';

        // Show loading indicator
        const loadingMessage = document.createElement('div');
        loadingMessage.classList.add('message', 'bot');
        loadingMessage.innerHTML = '<div class="message-content">Thinking...</div>';
        chatMessages.appendChild(loadingMessage);
        chatMessages.scrollTop = chatMessages.scrollHeight;

        try {
            // Debug: print the payload being sent
            console.log("Sending payload:", JSON.stringify({ messages }));
            // Send request to backend
            const response = await fetch('/api/chat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    ...(currentConversationID ? { 'X-Conversation-Id': currentConversationID } : {})
                },
                body: JSON.stringify({ messages }),
            });

            if (!response.ok) {
                throw new Error('Failed to get response from chatbot');
            }

            // Update conversation ID from response header if present
            const newConvId = response.headers.get('X-Conversation-Id');
            if (newConvId) {
                currentConversationID = newConvId;
            }

            const data = await response.json();
            
            // Remove loading indicator
            chatMessages.removeChild(loadingMessage);
            
            // Add bot response to chat
            addMessage(data.message);
            messages.push({ role: 'assistant', content: data.message });
        } catch (error) {
            // Remove loading indicator
            chatMessages.removeChild(loadingMessage);
            
            // Add error message
            addMessage('Sorry, there was an error communicating with the chatbot. Please try again later.');
            console.error('Error:', error);
        }
    }

    // Event listeners
    sendButton.addEventListener('click', () => {
        sendMessage(messageInput.value);
    });

    messageInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            sendMessage(messageInput.value);
        }
    });

    // Focus input field on load
    messageInput.focus();

    // --- Chat History Sidebar Logic ---
    let currentConversationID = null;

    const sidebar = document.getElementById('chat-sidebar');
    const sidebarToggle = document.getElementById('sidebar-toggle');
    const sidebarClose = document.getElementById('sidebar-close');
    const historySearchInput = document.getElementById('history-search-input');
    const historySearchButton = document.getElementById('history-search-button');
    const historyList = document.getElementById('history-list');

    sidebarToggle.addEventListener('click', () => {
        sidebar.style.display = 'block';
    });
    sidebarClose.addEventListener('click', () => {
        sidebar.style.display = 'none';
    });

    async function searchHistory(term) {
        const res = await fetch(`/api/history/search?q=${encodeURIComponent(term)}`);
        const data = await res.json();
        historyList.innerHTML = '';
        if (data.matches && data.matches.length > 0) {
            data.matches.forEach(id => {
                const li = document.createElement('li');
                li.textContent = id;
                li.style.cursor = 'pointer';
                li.onclick = () => loadHistory(id);
                historyList.appendChild(li);
            });
        } else {
            const li = document.createElement('li');
            li.textContent = 'No matches found.';
            historyList.appendChild(li);
        }
    }

    historySearchButton.addEventListener('click', () => {
        searchHistory(historySearchInput.value);
    });

    // Load and display a conversation's history
    async function loadHistory(conversationID) {
        const res = await fetch(`/api/history/${conversationID}`);
        const data = await res.json();
        chatMessages.innerHTML = '';
        if (data.history && data.history.length > 0) {
            data.history.forEach(line => {
                if (line.trim() === '') return;
                // Parse line: [timestamp] [role]: message
                const match = line.match(/^\[(.*?)\] \[(.*?)\]: (.*)$/);
                if (match) {
                    const role = match[2] === 'user' ? 'user' : 'bot';
                    addMessage(match[3], role === 'user');
                } else {
                    addMessage(line, false);
                }
            });
        }
        currentConversationID = conversationID;
        sidebar.style.display = 'none';
    }

    // Optionally, load all history on sidebar open
    sidebarToggle.addEventListener('click', () => {
        searchHistory('');
    });

    // Start new chat: clear conversation ID and chat window
    window.startNewChat = function() {
        currentConversationID = null;
        chatMessages.innerHTML = '';
        messages.length = 0;
        addMessage("Hello! I can help you manage your Cal.com calendar. Here are some things you can say:\n<ul><li>\"Help me book a meeting\"</li><li>\"Show me my scheduled events\"</li><li>\"Cancel my event at 3pm today\"</li><li>\"Reschedule my 2pm meeting to tomorrow at 4pm\"</li></ul>", false);
    }

    // --- End Chat History Sidebar Logic ---
}); 