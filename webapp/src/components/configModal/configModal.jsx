import React from 'react';
import PropTypes from 'prop-types';
import {
    Alert,
    Button,
    ControlLabel,
    FormControl,
    FormGroup,
    InputGroup,
    MenuItem,
    Modal,
    SplitButton,
    Tab,
    Tabs,
} from 'react-bootstrap';
import Constants from '../../constants';
import TimePicker from '../timePicker';
import request from 'superagent';
import style from './style.css';
import reactStyles from './style';
import SentryBoundary from '../../SentryBoundary';
import * as HttpStatus from 'http-status-codes';
import ToggleSwitch from '../toggleSwitch';
import Cookies from 'js-cookie';
import RRule from '../rRule';

const configModalCloseTimeout = 1000;
const timezones = require('../../../../timezones.json');

class ConfigModal extends (SentryBoundary, React.Component) {
    constructor(props) {
        super(props);
        this.state = this.getInitialState();

        // eslint-disable-next-line no-unused-vars
        const x = style;
    }

    static get REPORT_DISPLAY_NAMES() {
        return {
            user_aggregated: 'User Aggregated',
            type_aggregated: 'Type Aggregated',
        };
    }

    static get STATUS_DISPLAY_NAMES() {
        return {
            true: 'Enabled',
            false: 'Disabled',
        };
    }

    static get TIMEZONE_DISPLAY_NAMES() {
        const timezoneList = {};
        for (let i = 0; i < Object.keys(timezones).length; ++i) {
            timezoneList[timezones[i]['display_name']] = timezones[i]['value'];
        }
        timezoneList[''] = '-';
        return timezoneList;
    }

    getInitialState = () => {
        return {
            showSpinner: true,
            hasPermission: undefined,
            windowOpenTime: '00:00',
            windowCloseTime: '00:00',
            reportFormat: 'user_aggregated',
            sections: {},
            members: [],
            enabled: true,
            status: true,
            message: {
                show: false,
                text: '',
                type: 'info',
            },
            windowOpenReminderEnabled: true,
            windowCloseReminderEnabled: true,
            timezone: '',
            scheduleEnabled: false,
            schedule: '',
            rruleString: '',
            startDate: new Date().toISOString(),
        };
    };

    handleClose = () => {
        this.setState(this.getInitialState);
        this.props.close();
    };

    handleWindowOpenTimeChange = (time) => {
        this.setState({
            windowOpenTime: time,
        });
    };

    handleWindowCloseTimeChange = (time) => {
        this.setState({
            windowCloseTime: time,
        });
    };

    handleReportTypeChange = (reportType) => {
        this.setState({
            reportFormat: reportType,
        });
    };

    handleStatusChange = () => {
        this.setState({
            enabled: !this.state.enabled,
        });
    };

    handleTimezoneChange = (timezone) => {
        this.setState({timezone});
    };

    handleWindowCloseReminderChange = () => {
        this.setState({
            windowCloseReminderEnabled: !this.state.windowCloseReminderEnabled,
        });
    };

    handleWindowOpenReminderChange = () => {
        this.setState({
            windowOpenReminderEnabled: !this.state.windowOpenReminderEnabled,
        });
    };

    handleScheduleStatusChange = () => {
        this.setState({
            scheduleEnabled: !this.state.scheduleEnabled,
        });
    };

    handleRecurrenceChange = (rruleString, startDate) => {
        this.setState({
            rruleString,
            startDate,
        });
    };

    generateSections = (onChangeCallback) => {
        // eslint-disable-next-line no-shadow
        const style = reactStyles.getStyle();
        const sections = [];

        for (let i = 0; i <= Object.keys(this.state.sections).length; ++i) {
            sections.push(
                <FormGroup
                    key={i.toString()}
                    style={{...style.formGroup, ...style.sections}}
                >
                    <InputGroup>
                        <InputGroup.Addon>{(i + 1) + '.'}</InputGroup.Addon>
                        <FormControl
                            type={'text'}
                            name={`line${i + 1}`}
                            onChange={onChangeCallback}
                            value={this.state.sections[`line${i + 1}`] || ''}
                        />
                    </InputGroup>
                </FormGroup>,
            );
        }

        return sections;
    };

    handleSectionChange = (e) => {
        const sections = {...this.state.sections};
        sections[e.target.name] = e.target.value;
        this.setState({
            sections,
        });
    };

    componentDidUpdate(prevProp) {
        if (this.props.visible !== prevProp.visible && this.props.visible) {
            this.getStandupConfig()
                .then(() => {
                    this.setState({showSpinner: false});
                })
                .catch(() => {
                    this.setState({showSpinner: false});
                });
        }
    }

    getStandupConfig = () => {
        return new Promise((resolve) => {
            const url = `${this.props.siteURL}/${Constants.URL_STANDUP_CONFIG}?channel_id=${this.props.channelID}`;
            request
                .get(url)
                .withCredentials()
                .end((err, result) => {
                    if (result.ok) {
                        const standupConfig = result.body;
                        const state = {
                            hasPermission: true,
                            windowOpenTime: standupConfig.windowOpenTime,
                            windowCloseTime: standupConfig.windowCloseTime,
                            reportFormat: standupConfig.reportFormat,
                            members: standupConfig.members,
                            sections: {},
                            enabled: standupConfig.enabled,
                            status: standupConfig.enabled,
                            timezone: standupConfig.timezone,
                            windowOpenReminderEnabled: standupConfig.windowOpenReminderEnabled,
                            windowCloseReminderEnabled: standupConfig.windowCloseReminderEnabled,
                            scheduleEnabled: standupConfig.scheduleEnabled,
                            schedule: standupConfig.schedule,
                            rruleString: standupConfig.rruleString,
                            startDate: standupConfig.startDate,
                        };

                        for (let i = 0; i < standupConfig.sections.length; ++i) {
                            state.sections[`line${i + 1}`] = standupConfig.sections[i];
                        }

                        this.setState(state);
                    } else if (result.status === HttpStatus.NOT_FOUND) {
                        // fetch system default timezone
                        request
                            .get(`${this.props.siteURL}/${Constants.URL_GET_TIMEZONE}`)
                            .withCredentials()
                            .end((error, response) => {
                                if (response.ok) {
                                    const timezone = String(response.body);
                                    this.setState({
                                        timezone,
                                    });
                                } else if (error) {
                                    console.log(error);
                                }
                            });

                        this.setState({
                            hasPermission: true,
                        });
                    } else if (result.status === HttpStatus.UNAUTHORIZED) {
                        this.setState({
                            hasPermission: false,
                        });
                    }

                    resolve();
                });
        });
    };

    prepareStandupConfigPayload() {
        return {
            channelId: this.props.channelID,
            windowOpenTime: this.state.windowOpenTime,
            windowCloseTime: this.state.windowCloseTime,
            reportFormat: this.state.reportFormat,
            sections: Object.values(this.state.sections).map((x) => x.trim()).filter((x) => x !== ''),
            members: this.state.members,
            enabled: this.state.enabled,
            timezone: this.state.timezone,
            windowCloseReminderEnabled: this.state.windowCloseReminderEnabled,
            windowOpenReminderEnabled: this.state.windowOpenReminderEnabled,
            scheduleEnabled: this.state.scheduleEnabled,
            rruleString: this.state.rruleString,
            startDate: this.state.startDate,
        };
    }

    saveStandupConfig = (e) => {
        e.preventDefault();

        // hiding message section so animation can re-trigger on new message
        this.setState({
            message: {
                show: false,
            },
        });

        request
            .post(`${this.props.siteURL}/${Constants.URL_STANDUP_CONFIG}`)
            .withCredentials()
            .send(this.prepareStandupConfigPayload())
            .set('X-CSRF-Token', Cookies.get(Constants.MATTERMOST_CSRF_COOKIE))
            .set('Content-Type', 'application/json')
            .end((err, res) => {
                if (err) {
                    this.setState({
                        message: {
                            show: true,
                            text: 'An error occurred while saving standup config.\n' + err.response.text,
                            type: 'danger',
                        },
                    });
                } else {
                    this.setState({
                        message: {
                            show: true,
                            text: 'Standup config saved successfully!',
                            type: 'success',
                        },
                    });
                    setTimeout(this.handleClose, configModalCloseTimeout);
                }
            });
    };

    render() {
        // eslint-disable-next-line no-shadow
        const style = reactStyles.getStyle();
        const data = timezones.map((timezone) =>
            (
                <MenuItem
                    key={timezone.value}
                    eventKey={timezone.value}
                >
                    {timezone.display_name}
                </MenuItem>
            ),
        );

        let showStandupError = false;
        let standupErrorMessage = '';
        let standupErrorSubMessage = '';

        if (this.state.hasPermission === false) {
            showStandupError = true;
            standupErrorMessage = 'You do not have permission to perform this operation';
            standupErrorSubMessage = 'Only a channel admin can perform this operation';
        }

        const spinner =
            (<div style={style.spinner}>
                <img
                    src={`${this.props.siteURL}/${Constants.URL_SPINNER_ICON}`}
                    alt={'loading...'}
                />
            </div>);

        const errorMessage =
            (<span>
                <span style={style.standupErrorMessage}>{standupErrorMessage}</span>
                <br/><br/>
                <span>{standupErrorSubMessage}</span>
            </span>);

        return (
            <Modal
                show={this.props.visible}
                onHide={this.handleClose}
                backdrop={'static'}
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        {`${Constants.PLUGIN_DISPLAY_NAME} - Configure`}
                    </Modal.Title>
                </Modal.Header>

                <Modal.Body style={showStandupError ? {} : style.body}>
                    {/* in progress spinner */}
                    <span hidden={!this.state.showSpinner}>
                        {spinner}
                    </span>

                    {/* generic error message section */}
                    <span hidden={this.state.showSpinner || !showStandupError}>
                        {errorMessage}
                    </span>

                    <div hidden={this.state.showSpinner || !this.state.hasPermission || showStandupError}>
                        <Tabs id={'standup-config-tabs'}>
                            <Tab
                                eventKey={1}
                                title={'General'}
                            >
                                <FormGroup style={style.formGroup}>
                                    <ControlLabel style={style.controlLabel}>
                                        {'Enabled:'}
                                    </ControlLabel>
                                    <ToggleSwitch
                                        onChange={this.handleStatusChange}
                                        checked={this.state.enabled}
                                        theme={this.props.theme}
                                    />
                                </FormGroup>
                                <FormGroup style={style.formGroup}>
                                    <ControlLabel style={style.controlLabel}>
                                        {'Standup Schedule:'}
                                    </ControlLabel>
                                    <ToggleSwitch
                                        onChange={this.handleScheduleStatusChange}
                                        checked={this.state.scheduleEnabled}
                                        theme={this.props.theme}
                                    />
                                </FormGroup>
                                <FormGroup style={style.formGroup}>
                                    <ControlLabel style={style.controlLabel}>
                                        {'Standup Report Format:'}
                                    </ControlLabel>
                                    <SplitButton
                                        title={ConfigModal.REPORT_DISPLAY_NAMES[this.state.reportFormat]}
                                        onSelect={this.handleReportTypeChange}
                                        bsStyle={'link'}
                                    >
                                        <MenuItem eventKey={'user_aggregated'}>{'User Aggregated'}</MenuItem>
                                        <MenuItem eventKey={'type_aggregated'}>{'Type Aggregated'}</MenuItem>
                                    </SplitButton>
                                </FormGroup>
                                <FormGroup style={{...style.formGroup, ...style.formGroupNoMarginBottom}}>
                                    <ControlLabel style={style.controlLabel}>{'Sections:'}</ControlLabel>
                                </FormGroup>

                                <div style={style.sectionGroup}>
                                    {this.generateSections(this.handleSectionChange)}
                                </div>
                            </Tab>
                            <Tab
                                eventKey={2}
                                title={'Notifications'}
                            >
                                <FormGroup style={style.formGroup}>
                                    <ControlLabel style={style.controlLabel}>
                                        {'Window Open Reminder:'}
                                    </ControlLabel>
                                    <ToggleSwitch
                                        onChange={this.handleWindowOpenReminderChange}
                                        checked={this.state.windowOpenReminderEnabled}
                                        theme={this.props.theme}
                                    />
                                </FormGroup>
                                <FormGroup style={style.formGroup}>
                                    <ControlLabel style={style.controlLabel}>
                                        {'Window Close Reminder:'}
                                    </ControlLabel>
                                    <ToggleSwitch
                                        onChange={this.handleWindowCloseReminderChange}
                                        checked={this.state.windowCloseReminderEnabled}
                                        theme={this.props.theme}
                                    />
                                </FormGroup>
                            </Tab>
                            <Tab
                                eventKey={3}
                                title={'Schedule'}
                            >
                                <FormGroup style={style.formGroup}>
                                    <ControlLabel style={style.controlLabel}>
                                        {'Window Time:'}
                                    </ControlLabel>
                                    <TimePicker
                                        time={this.state.windowOpenTime}
                                        onChange={this.handleWindowOpenTimeChange}
                                        bsStyle={'link'}
                                    />
                                    <span style={style.controlLabelX}>{'to'}</span>
                                    <TimePicker
                                        time={this.state.windowCloseTime}
                                        onChange={this.handleWindowCloseTimeChange}
                                        bsStyle={'link'}
                                    />
                                </FormGroup>
                                <FormGroup style={style.formGroup}>
                                    <ControlLabel style={style.controlLabel}>
                                        {'Timezone:'}
                                    </ControlLabel>
                                    <SplitButton
                                        title={ConfigModal.TIMEZONE_DISPLAY_NAMES[this.state.timezone]}
                                        onSelect={this.handleTimezoneChange}
                                        bsStyle={'link'}
                                    >{data}
                                    </SplitButton>
                                </FormGroup>
                                <FormGroup>
                                    <RRule
                                        startDate={this.state.startDate}
                                        rrule={this.state.rruleString}
                                        onChange={this.handleRecurrenceChange}
                                    />
                                </FormGroup>
                            </Tab>
                        </Tabs>
                    </div>
                </Modal.Body>

                <Modal.Footer>
                    <Button
                        type='button'
                        onClick={this.handleClose}
                        variant={'primary'}
                    >
                        {'Cancel'}
                    </Button>
                    <Button
                        type='submit'
                        bsStyle='primary'
                        onClick={this.saveStandupConfig}
                    >
                        {'Save'}
                    </Button>
                </Modal.Footer>

                <Alert
                    bsStyle={this.state.message.type}
                    style={style.alert}
                    className={(this.state.message.show ? '' : 'hidden')}
                >
                    {this.state.message.text}
                </Alert>
            </Modal>
        );
    }
}

ConfigModal.propTypes = {
    channelID: PropTypes.string.isRequired,
    currentUserId: PropTypes.string.isRequired,
    close: PropTypes.func.isRequired,
    visible: PropTypes.bool,
    siteURL: PropTypes.string.isRequired,
};

export default ConfigModal;
