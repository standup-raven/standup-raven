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
} from 'react-bootstrap';
import Constants from '../../constants';
import TimePicker from '../timePicker';
import request from 'superagent';
import style from './style.css';
import reactStyles from './style';
import SentryBoundary from '../../SentryBoundary';
import * as HttpStatus from 'http-status-codes';
import Cookies from 'js-cookie';

const configModalCloseTimeout = 1000;

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

    getInitialState = () => {
        return {
            showSpinner: true,
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

    handleStatusChange = (status) => {
        this.setState({
            enabled: status,
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
            const url = `${Constants.URL_STANDUP_CONFIG}?channel_id=${this.props.channelID}`;
            request
                .get(url)
                .withCredentials()
                .end((err, result) => {
                    if (result.ok) {
                        const standupConfig = result.body;
                        const state = {
                            windowOpenTime: standupConfig.windowOpenTime,
                            windowCloseTime: standupConfig.windowCloseTime,
                            reportFormat: standupConfig.reportFormat,
                            members: standupConfig.members,
                            sections: {},
                            enabled: standupConfig.enabled,
                            status: standupConfig.enabled,
                        };

                        for (let i = 0; i < standupConfig.sections.length; ++i) {
                            state.sections[`line${i + 1}`] = standupConfig.sections[i];
                        }

                        this.setState(state);
                    } else if (result.status !== HttpStatus.NOT_FOUND) {
                        console.log(err);
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
            .post(Constants.URL_STANDUP_CONFIG)
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

        const showStandupError = false;
        const standupErrorMessage = '';
        const standupErrorSubMessage = '';

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
                    <div
                        className={this.state.showSpinner ? '' : 'hidden'}
                        style={style.spinner}
                    >
                        <img
                            src={Constants.URL_SPINNER_ICON}
                            alt={'loading...'}
                        />
                    </div>

                    <div className={this.state.showSpinner ? 'hidden' : ''}>
                        <span className={showStandupError ? '' : 'hidden'}>
                            <span style={style.standupErrorMessage}>{standupErrorMessage}</span>
                            <br/><br/>
                            <span>{standupErrorSubMessage}</span>
                        </span>

                        <span className={showStandupError ? 'hidden' : ''}>
                            <FormGroup style={style.formGroup}>
                                <ControlLabel style={style.controlLabel}>
                                    {'Status:'}
                                </ControlLabel>

                                <SplitButton
                                    title={ConfigModal.STATUS_DISPLAY_NAMES[this.state.enabled]}
                                    onSelect={this.handleStatusChange}
                                    bsStyle={'link'}
                                >
                                    <MenuItem eventKey={true}>{ConfigModal.STATUS_DISPLAY_NAMES[true]}</MenuItem>
                                    <MenuItem eventKey={false}>{ConfigModal.STATUS_DISPLAY_NAMES[false]}</MenuItem>
                                </SplitButton>
                            </FormGroup>

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
                        </span>
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
};

export default ConfigModal;
