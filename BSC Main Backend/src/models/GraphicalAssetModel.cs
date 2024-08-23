using System.Reflection.Metadata;

namespace BSC_Main_Backend.Models
{
    public class GraphicalAssetModel
    {
        public int Id { get; set; } // Maps to SERIAL PRIMARY KEY
        public string Alias { get; set; } // Maps to alias VARCHAR(255) NOT NULL
        public string Type { get; set; } // Maps to type VARCHAR(255) NOT NULL
        public Blob Blob { get; set; } // Maps to blob BYTEA
        public string UseCase { get; set; } // Maps to useCase VARCHAR(255)
        public int Width { get; set; } // Maps to width INT NOT NULL
        public int Height { get; set; } // Maps to height INT NOT NULL
        public bool HasLOD { get; set; } // Maps to hasLOD BOOLEAN NOT NULL
    }
}