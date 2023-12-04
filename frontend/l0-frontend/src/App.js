import './App.css'
import {Component} from "react"
import { JsonView, allExpanded, darkStyles, defaultStyles } from 'react-json-view-lite'
import 'react-json-view-lite/dist/index.css'
import Button from 'react-bootstrap/Button'
import 'bootstrap/dist/css/bootstrap.min.css'
import Modal from 'react-bootstrap/Modal';
import Form from 'react-bootstrap/Form';
import OrderModal from "./OrderModal";

class App extends Component {
  state = {
    orders: [],
    showModal: false,
    orderId: '',
    actionTypeText: 'Get all orders'
  }

  async loadAllOrders() {
    const getAllOrdersAddress = "/order/"
    const response = await fetch(getAllOrdersAddress)
    const body = await response.json()
    this.setState({orders: body, actionTypeText: `Get all orders: ${body.length} elements`})
  }

  async loadCachedOrders() {
    const getCachedOrdersAddress = "/order/entire-cache"
    const response = await fetch(getCachedOrdersAddress)
    const body = await response.json()
    this.setState({orders: body, actionTypeText: `Get all cached orders: ${body.length} elements`})
  }

  async getOrderById(id) {
    const getOrderAddress = `/order/${id}`
    const response = await fetch(getOrderAddress)
    if (response.status === 404 || response.ok === false) {
      this.setState({orders: [], actionTypeText: 'Cannot find order with this id'})
      return
    }
    const body = await response.json()
    this.setState({orders: [body], actionTypeText: `Get order by id`})
  }

  handleShowModal = () => {
    this.setState({ showModal: true })
  }

  handleCloseModal = () => {
    this.setState({ showModal: false })
  }

  handleIdChange = (event) => {
    this.setState({ orderId: event.target.value })
  }

  handleGetOrderById = () => {
    const { orderId } = this.state
    if (orderId) {
      this.getOrderById(orderId)
      this.handleCloseModal()
    }
  }

  render() {
    const {orders, showModal, orderId, actionTypeText} = this.state
    return (
        <div className="App">
          <header className="App-header">
            <h2 style={{marginTop: '20px'}}>Choose your action</h2>
            <div style={{}}>
              <Button variant="primary" onClick={() => this.loadAllOrders()}>Get all Orders</Button>{' '}
              <Button variant="primary" style={{marginLeft: '6px'}} onClick={() => this.loadCachedOrders()}
              >Get all cached Orders</Button>{' '}
              <Button variant="primary" style={{marginLeft: '6px'}} onClick={this.handleShowModal}
              >Get order by id</Button>{' '}
            </div>
            <div className="App-intro" style={{marginTop: '20px'}}>
              <h2 style={{marginTop: '20px'}}>{actionTypeText}</h2>
              {orders.map((order, index) =>
                  <div style={{fontSize: '16px'}}>
                    <h3 style={{marginTop: '20px'}}>This is order {index + 1}</h3>
                    <div style={{textAlign: 'left', marginTop: '20px'}}>
                      <JsonView data={order} shouldExpandNode={allExpanded} style={defaultStyles}/>
                    </div>
                  </div>
              )}
            </div>

            <OrderModal
                showModal={showModal}
                handleCloseModal={this.handleCloseModal}
                orderId={orderId}
                handleIdChange={this.handleIdChange}
                handleGetOrderById={this.handleGetOrderById}
            />
          </header>
        </div>
    )
  }
}

export default App
