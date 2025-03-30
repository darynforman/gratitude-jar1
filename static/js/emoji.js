// Emoji selection functionality
function selectEmoji(button) {
    console.log('Emoji button clicked:', button.dataset.emoji);
    
    // Remove selected class from all emoji buttons
    document.querySelectorAll('.emoji-btn').forEach(btn => {
        btn.classList.remove('bg-indigo-100', 'ring-2', 'ring-indigo-500');
    });
    
    // Add selected class to clicked button
    button.classList.add('bg-indigo-100', 'ring-2', 'ring-indigo-500');
    
    // Update hidden input with selected emoji
    document.getElementById('selected-emoji').value = button.dataset.emoji;
} 