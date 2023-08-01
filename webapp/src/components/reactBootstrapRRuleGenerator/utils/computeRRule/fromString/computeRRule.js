import {rrulestr as RRuleObjectFromString} from 'rrule';
import moment from 'moment';

import {DATE_TIME_FORMAT} from '../../../constants/index';
import computeStartOnDate from './computeStartOnDate';
import computeFrequency from './computeFrequency';
import computeYearlyMode from './computeYearlyMode';
import computeYearlyOnMonth from './computeYearlyOnMonth';
import computeYearlyOnMonthday from './computeYearlyOnMonthday';
import computeYearlyOnTheMonth from './computeYearlyOnTheMonth';
import computeYearlyOnTheMonthday from './computeYearlyOnTheMonthday';
import computeYearlyOnTheWhich from './computeYearlyOnTheWhich';
import computeMonthlyMode from './computeMonthlyMode';
import computeMonthlyInterval from './computeMonthlyInterval';
import computeMonthlyOnDay from './computeMonthlyOnDay';
import computeMonthlyOnTheDay from './computeMonthlyOnTheDay';
import computeMonthlyOnTheWhich from './computeMonthlyOnTheWhich';
import computeWeeklyInterval from './computeWeeklyInterval';
import computeWeeklyDays from './computeWeeklyDays';
import computeWeekStartDay from './computeWeekStartDay';
import computeDailyInterval from './computeDailyInterval';
import computeHourlyInterval from './computeHourlyInterval';
import computeEndMode from './computeEndMode';
import computeEndAfter from './computeEndAfter';
import computeEndOnDate from './computeEndOnDate';

const computeRRule = (data, rrule) => {
    if (!rrule) {
        return data;
    }

    let newDataObj;
    try {
        const rruleObj = RRuleObjectFromString(rrule).origOptions;

        newDataObj = {
            ...data,
            start: {
                ...data.start,
                onDate: {
                    date: moment(computeStartOnDate(data, rruleObj)).format(DATE_TIME_FORMAT),
                    options: {
                        ...data.start.onDate.options,
                        weekStartsOnSunday: computeWeekStartDay(data, rruleObj),
                    },
                },
            },
            repeat: {
                ...data.repeat,
                frequency: computeFrequency(data, rruleObj),
                yearly: {
                    ...data.repeat.yearly,
                    mode: computeYearlyMode(data, rruleObj),
                    on: {
                        month: computeYearlyOnMonth(data, rruleObj),
                        day: computeYearlyOnMonthday(data, rruleObj),
                    },
                    onThe: {
                        month: computeYearlyOnTheMonth(data, rruleObj),
                        day: computeYearlyOnTheMonthday(data, rruleObj),
                        which: computeYearlyOnTheWhich(data, rruleObj),
                    },
                },
                monthly: {
                    ...data.repeat.monthly,
                    mode: computeMonthlyMode(data, rruleObj),
                    interval: computeMonthlyInterval(data, rruleObj),
                    on: {
                        day: computeMonthlyOnDay(data, rruleObj),
                    },
                    onThe: {
                        day: computeMonthlyOnTheDay(data, rruleObj),
                        which: computeMonthlyOnTheWhich(data, rruleObj),
                    },
                },
                weekly: {
                    interval: computeWeeklyInterval(data, rruleObj),
                    days: computeWeeklyDays(data, rruleObj),
                    options: {
                        weekStartsOnSunday: computeWeekStartDay(data, rruleObj),
                    },
                },
                daily: {
                    interval: computeDailyInterval(data, rruleObj),
                },
                hourly: {
                    interval: computeHourlyInterval(data, rruleObj),
                },
            },
            end: {
                ...data.end,
                mode: computeEndMode(data, rruleObj),
                after: computeEndAfter(data, rruleObj),
                onDate: {
                    date: moment(computeEndOnDate(data, rruleObj)).format(DATE_TIME_FORMAT),
                    options: {
                        ...data.end.onDate.options,
                        weekStartsOnSunday: computeWeekStartDay(data, rruleObj),
                    },
                },
            },
            options: {
                ...data.options,
                weekStartsOnSunday: computeWeekStartDay(data, rruleObj),
            },
            error: null,
        };
    } catch (e) {
        return {...data, error: {value: rrule, message: e}};
    }

    return newDataObj;
};

export default computeRRule;
