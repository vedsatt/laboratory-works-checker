#include <boost/asio.hpp>
#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/json_parser.hpp>
#include <iostream>
#include <sstream>

using namespace boost::asio;
using ip::tcp;
namespace pt = boost::property_tree;

struct LabWork {
    int id;
    std::string code;
    int task1, task2, task3;
};

struct CheckerResponse {
    int id;
    std::string status;
    std::string msg;
};

std::string create_response_json(const CheckerResponse& cr) {
    pt::ptree root;
    root.put("id", cr.id);        // Число, а не строка!
    root.put("status", cr.status);
    root.put("msg", cr.msg);
    
    std::ostringstream oss;
    pt::write_json(oss, root, false); // false - убираем pretty-formatting
    return oss.str();
}

int main() {
    try {
        io_context io;
        tcp::acceptor acceptor(io, tcp::endpoint(tcp::v4(), 8080));
        std::cout << "Server started on port 8080\n";

        while (true) {
            tcp::socket socket(io);
            acceptor.accept(socket);
            
            try {
                // Чтение запроса
                streambuf buf;
                read_until(socket, buf, '\n');
                
                std::istream is(&buf);
                std::string request;
                std::getline(is, request);
                
                std::cout << "Received request: " << request << std::endl;
                
                // Формируем ответ (в реальном коде здесь должен быть парсинг)
                CheckerResponse cr{
                    123,           // ID (число!)
                    "success",     // Статус
                    "All checks passed" // Сообщение
                };
                
                // Создаем JSON
                std::string response = create_response_json(cr);
                response += "\n";  // Добавляем разделитель
                
                // Отправляем
                write(socket, buffer(response));
                std::cout << "Sent response: " << response;
                
                // Закрываем соединение
                socket.shutdown(tcp::socket::shutdown_send);
            } catch (const std::exception& e) {
                std::cerr << "Error: " << e.what() << std::endl;
            }
        }
    } catch (const std::exception& e) {
        std::cerr << "Server error: " << e.what() << std::endl;
        return 1;
    }
    return 0;
}