// JS here

$( document ).ready(function() {

    // Submit user selected answer
    $('.submit-answer').click(function(e) {
        e.preventDefault();
        $('#user_answer').val($(this).data('value'));
        $('#next_question').val($(this).data('next'));
        $('#questionForm').submit();
    });

});