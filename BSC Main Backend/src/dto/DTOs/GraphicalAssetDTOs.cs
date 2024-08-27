using System.Reflection.Metadata;

namespace BSC_Main_Backend.Dtos
{
    public class GraphicalAssetDTOs
    {
        public int Id { get; set; }
        public string Alias { get; set; }
        public string Type { get; set; }
        public string UseCase { get; set; }
        public int Width { get; set; }
        public int Height { get; set; }
        public bool HasLOD { get; set; }
    }
    
    public class GraphicalAssetCreateDTO
    {
        public string Alias { get; set; }
        public string Type { get; set; }
        public Blob Blob { get; set; }
        public string UseCase { get; set; }
        public int Width { get; set; }
        public int Height { get; set; }
        public bool HasLOD { get; set; }
    }
    
    public class GraphicalAssetUpdateDTO
    {
        public string Alias { get; set; }
        public string Type { get; set; }
        public Blob Blob { get; set; }
        public string UseCase { get; set; }
        public int Width { get; set; }
        public int Height { get; set; }
        public bool HasLOD { get; set; }
    }
}