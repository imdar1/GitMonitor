# GitMonitor Design Docs

## Software Requirements Docs
### Use case diagram
![Use Case Diagram](res/use_case.svg)
There are 6 main use-cases with its corresponding sequence diagram
- [Atur direktori proyek Git](#use-case-atur-direktori-proyek-git)
- [Atur branch proyek Git](#use-case-atur-branch-proyek-git)
- [Atur task proyek](#use-case-atur-task-proyek)
- [Lihat timeline task proyek](#use-case-lihat-timeline-task-proyek)
- [Lihat informasi proyek berdasarkan data Git](#use-case-lihat-informasi-proyek-berdasarkan-data-git)
- [Atur periode pembaruan data Git](#use-case-lihat-informasi-proyek-berdasarkan-data-git)

### Sequence diagram
#### Use case "Atur direktori proyek Git"
![Sequence Diagram atur direktori proyek Git](res/sequence_diagram1.svg)

#### Use case "Atur branch proyek Git"
![Sequence Diagram atur branch proyek Git](res/sequence_diagram2.svg)

#### Use case "Atur task proyek"
![Sequence Diagram atur task proyek](res/sequence_diagram3.svg)

#### Use case "Lihat timeline task proyek"
![Sequence Diagram lihat timeline task proyek](res/sequence_diagram4.svg)

#### Use case "Lihat informasi proyek berdasarkan data Git"
![Sequence Diagram lihat informasi proyek berdasarkan data Git](res/sequence_diagram5.svg)

#### Use case "Atur periode pembaruan data Git"
![Sequence Diagram atur periode pembaruan data Git](res/sequence_diagram6.svg)

## Technical Docs
### High-level architecture
![High-level architecture](res/GitMonitor_architecture.png)

Each layer in application is tightly coupled. Users can only interact with the application using application UI provided. The application logics will serve all the logics needed to run the application, performing Git operations to a Git repositoriy and read/write operations to the embedded database.

### ER Diagram
![ER Diagram](res/er_diagram.svg)
