#include "server.hpp"
#include <boost/asio.hpp>
#include <boost/bind.hpp>
#include <iostream>
#include <sstream>
#include <json/json.h>

Server::Server(boost::asio::io_context& io_context, short port)
    : acceptor_(io_context, tcp::endpoint(tcp::v4(), port)) {
    start_accept();
}

void Server::start_accept() {
    tcp::socket socket(acceptor_.get_executor());
    acceptor_.async_accept(socket,
        boost::bind(&Server::handle_accept, this, std::move(socket),
        boost::asio::placeholders::error));
}

void Server::handle_accept(tcp::socket socket, const boost::system::error_code& error) {
    if (!error) {
        auto buffer = boost::asio::buffer(new char[1024], 1024);
        socket.async_read_some(buffer,
            boost::bind(&Server::handle_read, this, std::ref(socket),
                boost::asio::placeholders::error,
                boost::asio::placeholders::bytes_transferred));
    }
    start_accept();
}

void Server::handle_read(tcp::socket& socket, const boost::system::error_code& error, size_t bytes_transferred) {
    if (!error && bytes_transferred > 0) {
        std::string data(boost::asio::buffer_cast<const char*>(socket.receive_buffer().data()), bytes_transferred);
        
        Json::Value root;
        Json::CharReaderBuilder builder;
        std::string errors;
        std::istringstream iss(data);
        
        if (Json::parseFromStream(builder, iss, &root, &errors)) {
            LabWork labWork;
            labWork.ID = root["id"].asInt();
            labWork.Code = root["code"].asString();
            labWork.Task1 = root["task1"].asInt();
            labWork.Task2 = root["task2"].asInt();
            labWork.Task3 = root["task3"].asInt();
            
            process_lab_work(labWork, socket);
        }
    }
}

void Server::process_lab_work(const LabWork& labWork, tcp::socket& socket) {
    CheckerResponse response;
    response.ID = labWork.ID;
    
    // Простая логика проверки - все задания должны быть больше 0
    if (labWork.Task1 > 0 && labWork.Task2 > 0 && labWork.Task3 > 0) {
        response.Status = "success";
        response.Msg = "All tasks completed successfully";
    } else {
        response.Status = "error";
        response.Msg = "Some tasks are not completed";
    }
    
    send_response(response, socket);
}

void Server::send_response(const CheckerResponse& response, tcp::socket& socket) {
    Json::Value root;
    root["id"] = response.ID;
    root["status"] = response.Status;
    root["msg"] = response.Msg;
    
    Json::StreamWriterBuilder writer;
    std::string json_str = Json::writeString(writer, root);
    
    boost::asio::write(socket, boost::asio::buffer(json_str));
}