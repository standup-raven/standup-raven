import {MONTHS} from '../../../constants/index';

const computeYearlyOnThe = (onThe) => {
    const repeat = {};

    switch (onThe.which) {
        case 'First':
            repeat.bysetpos = 1;
            break;
        case 'Second':
            repeat.bysetpos = 2;
            break;
        case 'Third':
            repeat.bysetpos = 3;
            break;
        case 'Fourth':
            repeat.bysetpos = 4;
            break;
        case 'Last':
            repeat.bysetpos = -1;
            break;
        default:
            break;
    }

    switch (onThe.day) {
        case 'Monday':
            repeat.byweekday = [0];
            break;
        case 'Tuesday':
            repeat.byweekday = [1];
            break;
        case 'Wednesday':
            repeat.byweekday = [2];
            break;
        case 'Thursday':
            repeat.byweekday = [3];
            break;
        case 'Friday':
            repeat.byweekday = [4];
            break;
        case 'Saturday':
            repeat.byweekday = [5];
            break;
        case 'Sunday':
            repeat.byweekday = [6];
            break;
        case 'Day':
            repeat.byweekday = [0, 1, 2, 3, 4, 5, 6];
            break;
        case 'Weekday':
            repeat.byweekday = [0, 1, 2, 3, 4];
            break;
        case 'Weekend day':
            repeat.byweekday = [5, 6];
            break;
        default:
            break;
    }

    repeat.bymonth = MONTHS.indexOf(onThe.month) + 1;

    return repeat;
};

export default computeYearlyOnThe;
