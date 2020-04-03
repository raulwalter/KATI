// JS here

$(document).ready(function () {

    // Submit user selected answer
    $('.submit-answer').click(function (e) {
        e.preventDefault();
        $('#user_answer').val($(this).data('value'));
        $('#next_question').val($(this).data('next'));
        $('#questionForm').submit();
    });

    $('.qsubmit').click(function (e) {
        e.preventDefault();
        var id = parseInt($(e.currentTarget).attr('id'));
        if (id == 1) {
            var code = $('#user_answer').val();
            if(validate(code)) {
                $('#questionForm').submit();
            } else {
                $('#user_answer').addClass('is-invalid');
            }
        } else {
            $('#questionForm').submit();
        }
    });

});

/**
     *  Validate estonian id code
     *
     *  @return boolean
     */
validate = function (code) {
    if (code.length !== 11) {
        return false;
    }
    var control = this.getControlNumber(code);
    if (control !== parseInt(code.charAt(10))) {
        return false;
    }

    var year = Number(code.substr(1, 2));
    var month = Number(code.substr(3, 2));
    var day = Number(code.substr(5, 2));
    var birthDate = getBirthday(code);
    return year === birthDate.getFullYear() % 100 && birthDate.getMonth() + 1 === month && day === birthDate.getDate();
};

/**
   *  Get gender from code
   *
   *  @return string
   */
getGender = function (code) {
    var firstNumber = code.charAt(0),
        result = '';
    switch (firstNumber) {
        case '1':
        case '3':
        case '5':
            result = 'male';
            break;
        case '2':
        case '4':
        case '6':
            result = 'female';
            break;
        default:
            result = 'unknown';
    }
    return result;
};

/**
 *  Get the age in years.
 *
 *  @return number
 */
getAge = function (code) {
    return Math.floor((new Date().getTime() - getBirthday(code).getTime()) / (86400 * 1000) / 365.25);
};

/**
 *  Get birthday from code.
 *
 *  @return Date
 */
getBirthday = function (code) {
    var year = parseInt(code.substring(1, 3)),
        month = parseInt(code.substring(3, 5).replace(/^0/, '')) - 1,
        day = code.substring(5, 7).replace(/^0/, ''),
        firstNumber = code.charAt(0);

    if (firstNumber === '1' || firstNumber === '2') {
        year += 1800;
    } else if (firstNumber === '3' || firstNumber === '4') {
        year += 1900;
    } else if (firstNumber === '5' || firstNumber === '6') {
        year += 2000;
    } else {
        year += 2100;
    }

    return new Date(year, month, day);
};

/**
     *  Gets the control number from code.
     *
     *  @return number
     */
this.getControlNumber = function (code) {
    var multiplier1 = [1, 2, 3, 4, 5, 6, 7, 8, 9, 1],
        multiplier2 = [3, 4, 5, 6, 7, 8, 9, 1, 2, 3],
        mod,
        total = 0;

    for (var i = 0; i < 10; i++) {
        total += code.charAt(i) * multiplier1[i];
    }
    mod = total % 11;

    total = 0;
    if (10 === mod) {
        for (i = 0; i < 10; i++) {
            total += code.charAt(i) * multiplier2[i];
        }
        mod = total % 11;
        if (10 === mod) {
            mod = 0;
        }
    }
    return mod;
};