import React, { Component } from 'react';
import PropTypes from 'prop-types';

class User extends Component {
    render() {
        return (
            <li>
                <a>{this.props.user.name}</a>
            </li>
        );
    }
}

User.propTypes = {
    user: PropTypes.object.isRequired
}

export default User