namespace ASP.NETCoreWebAPI.Models
{
    public class Achievement
    {
        public int Id { get; set; } // Primary Key
        public required string Title { get; set; } // Attribute - Required, title of the achievement
        public string? Description { get; set; } // Attribute - Nullable, description of the achievement
        public GraphicalAsset? Icon { get; set; } // Foreign Key (GraphicalAsset) - Nullable, URL or path to the icon representing the achievement
        public bool IsTutorialCompleted { get; set; } // Attribute - Indicates if this achievement is for completing the tutorial
    }
}
