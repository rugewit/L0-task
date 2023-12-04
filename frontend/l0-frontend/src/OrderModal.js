// OrderModal.js
import React from 'react';
import Modal from 'react-bootstrap/Modal';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';

class OrderModal extends React.Component {
    render() {
        const { showModal, handleCloseModal, orderId, handleIdChange, handleGetOrderById } = this.props;

        return (
            <Modal show={showModal} onHide={handleCloseModal}>
                <Modal.Header closeButton>
                    <Modal.Title>Enter Order ID</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <Form>
                        <Form.Group controlId="formOrderId">
                            <Form.Label>Order ID:</Form.Label>
                            <Form.Control
                                type="text"
                                placeholder="Enter Order ID"
                                value={orderId}
                                onChange={handleIdChange}
                            />
                        </Form.Group>
                    </Form>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="secondary" onClick={handleCloseModal}>
                        Close
                    </Button>
                    <Button variant="primary" onClick={handleGetOrderById}>
                        Get Order
                    </Button>
                </Modal.Footer>
            </Modal>
        );
    }
}

export default OrderModal;
