const computeYearlyOnTheMonthday = (data, rruleObj) => {
    if (rruleObj.freq !== 0 || !rruleObj.byweekday) {
        return data.repeat.yearly.onThe.day;
    }

    const weekdays = rruleObj.byweekday.map((weekday) => weekday.weekday).join(',');

    switch (weekdays) {
        case '0': {
            return 'Monday';
        }
        case '1': {
            return 'Tuesday';
        }
        case '2': {
            return 'Wednesday';
        }
        case '3': {
            return 'Thursday';
        }
        case '4': {
            return 'Friday';
        }
        case '5': {
            return 'Saturday';
        }
        case '6': {
            return 'Sunday';
        }
        case '0,1,2,3,4,5,6': {
            return 'Day';
        }
        case '0,1,2,3,4': {
            return 'Weekday';
        }
        case '5,6': {
            return 'Weekend day';
        }
        default: {
            return data.repeat.yearly.onThe.day;
        }
    }
};

export default computeYearlyOnTheMonthday;
